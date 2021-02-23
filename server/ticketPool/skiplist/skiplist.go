// @Author: KongLingWen
// @Created at 2021/2/10
// @Modified at 2021/2/10

package skiplist

import (
	"fmt"
	"math/rand"
	"ticketPool/persistence"
)

type SkipList struct {
	head 			*Index
	maxLevel 		int
	requestChan 	chan *Request
	trainId			uint32
	date			string
	seatTypeId		uint32
}

type Node struct {
	Key   uint64
	Value []string
	Next  *Node
}


type Index struct {
	node 	*Node
	down 	*Index
	right 	*Index
}


type Request struct {
	option		string
	args		[]interface{}
}


func NewSkipList(trainId uint32, date string, seatTypeId uint32) *SkipList{
	head := &Index{
		node:  &Node{},
		down:  nil,
		right: nil,
	}
	sl := &SkipList{
		head:        head,
		maxLevel:    1,
		requestChan: make(chan *Request, 100),
		trainId:    trainId,
		date:       date,
		seatTypeId: seatTypeId,
	}
	sl.DealWithRequest()
	return sl
}

func(sl *SkipList) Do(cmd string, args ...interface{}){
	sl.requestChan <- &Request{
		option: cmd,
		args:   args,
	}
}

func(sl *SkipList) DealWithRequest(){
	go func() {
		for {
			req := <-sl.requestChan
			switch req.option {
			case "Put":
				key := req.args[0].(uint64)
				value := req.args[1].([]string)
				if key == 0 {
					continue
				}
				sl.Put(key, value)
				persistence.Do(&persistence.PersistentRequest{
					Option:     "INSERT",
					TrainId:    sl.trainId,
					Date:       sl.date,
					SeatTypeId: sl.seatTypeId,
					Key:        key,
					Value:      value,
				})
			case "Allocate":
				key := req.args[0].(uint64)
				count := req.args[1].(int)
				respChan := req.args[2].(chan *Node)
				node := sl.Allocate(key, count)
				respChan <- node
				for ; node != nil; node = node.Next{
					persistence.Do(&persistence.PersistentRequest{
						Option:     "DELETE",
						TrainId:    sl.trainId,
						Date:       sl.date,
						SeatTypeId: sl.seatTypeId,
						Key:        node.Key,
						Value:      node.Value,
					})
				}
			case "Search":
				key := req.args[0].(uint64)
				respChan := req.args[1].(chan int32)
				respChan <- sl.Search(key)
			}
		}
	}()
}

func(sl *SkipList) Get(key uint64) []string {
	b := sl.findPredecessor(key)
	n := b.Next

	for {
		if n == nil {
			return nil
		}
		c := cpr(key, n.Key)
		if c == 0 {
			return n.Value
		}
		if c < 0 {
			return nil
		}
		b = n
		n = n.Next
	}
}


func(sl *SkipList) Put(key uint64, value []string) {
	var node *Node

	b := sl.findPredecessor(key)
	n := b.Next

	for {
		if n != nil {
			c := cpr(key, n.Key)

			if c > 0 {
				b = n
				n = n.Next
				continue
			}

			if c == 0 {
				n.Value = append(n.Value, value...)
				return
			}
		}

		node = &Node{
			Key:   key,
			Value: value,
			Next:  n,
		}
		b.Next = node
		break
	}

	// 新插入节点时随机决定是否加入索引节点
	rnd := rand.Uint32()
	if rnd & 0x80000001 == 0 {
		level := 1
		for rnd >>= 1; (rnd & 1) != 0; rnd>>=1 {
			level++
		}
		var idx *Index

		// 随机level <= 当前 level 时，生成随机 level 个索引节点
		if level <= sl.maxLevel {
			for i := 0; i < level; i++{
				idx = &Index{
					node:  node,
					down:  idx,
					right: nil,
				}
			}
		} else {
			// 随机 level 大于当前 level 时
			// 增加索引层数
			level = sl.maxLevel + 1
			// 生成当 level 个索引节点
			for i := 0; i < level; i++ {
				idx = &Index{
					node:  node,
					down:  idx,
					right: nil,
				}
			}
			// 增加原头索引高度
			h := sl.head

			newH := &Index{
				node:  h.node,
				down:  h,
				right: nil,
			}
			sl.head = newH
			sl.maxLevel++
		}

		// 寻找索引插入点
		// 从头结点开始（头结点 level 为最高 level），每层向右寻找插入位置，找到后插入，然后进入下一层
		q := sl.head
		r := q.right
		hLevel := sl.maxLevel
		for {

			if r != nil {
				n := r.node
				c := cpr(key, n.Key)
				if c > 0 {
					q = r
					r = r.right
					continue
				}
			}
			if hLevel == level {
				idx.right = r
				q.right = idx
				idx = idx.down
				if level--; level < 1 {
					break
				}
			}
			hLevel--

			q = q.down
			r = q.right
		}
	}



}

func(sl *SkipList) Remove (key uint64) []string {
	b := sl.findPredecessor(key)
	n := b.Next

	for{
		if n == nil {
			return nil
		}
		c := cpr(key, n.Key)
		if c < 0 {
			return nil
		}
		if c > 0 {
			b = n
			n = n.Next
			continue
		}

		v := n.Value

		n.Value = nil
		n.appendMarker(n.Next)
		b.Next = n.Next

		sl.findPredecessor(key) // clean Index
		if sl.head.right == nil {
			sl.tryReduceLevel()
		}
		return v
	}
}

func(sl *SkipList) tryReduceLevel() {
	if sl.maxLevel <= 3 {
		return
	}
	d := sl.head.down

	if d == nil && d.down == nil {
		sl.head = d
	}
}

// return Node  -->  Node.Key <= Key < Node.Next.Key
func(sl *SkipList) findPredecessor(key uint64) *Node {
	var q, r, d *Index
	q = sl.head
	r = q.right
	for {
		if r != nil {
			n := r.node
			k := n.Key
			if n.Value == nil { // clean index
				q.right = r.right
				r = q.right
				continue
			}
			if cpr(key,k) > 0 { // 与 q 的右节点的 Key 比较，Key > k 则向右寻找
				q = r
				r = r.right
				continue
			}
		}
		// r == nil 或 Key <= k    r = q.right n == r.node
		if d = q.down; d == nil { // 进入下一层
			return q.node
		}
		q = d
		r = d.right
	}
}

func cpr(x, y uint64) int {
	if x < y {
		return -1
	} else if x > y {
		return 1
	}
	return 0
}


func(sl *SkipList) Allocate(key uint64, count int) *Node {
	b := sl.findPredecessor(key)
	n := b.Next

	var node *Node

	for{
		if n == nil {
			return node
		}
		c := cpr(key, n.Key)
		if c <= 0 {
			if key & n.Key == key {
				v := n.Value
				lenV := len(n.Value)
				if lenV > count {
					// 当前节点的座位数量满足需求，分配并返回
					retV := make([]string, count)
					copy(retV,v[lenV-count:])
					n.Value = v[:lenV-count]

					node = &Node{
						Key:   n.Key,
						Value: retV,
						Next:  node,
					}
					return node

				} else {
					// lenV <= count, 分配 lenV 个座位，删除节点，继续寻找
					node = &Node{
						Key:   n.Key,
						Value: n.Value,
						Next:  node,
					}
					count -= lenV

					// value 置为 nil，添加标记，用于清除索引节点
					n.Value = nil
					f := n.Next
					n.appendMarker(f)
					b.Next = f
					n = f

					sl.findPredecessor(key) // clean Index
					if sl.head.right == nil {
						sl.tryReduceLevel()
					}

					if count == 0{
						return node
					}
					continue
				}
			}
		}
		b = n
		n = n.Next
	}

}


func(sl *SkipList) Search(key uint64) int32 {
	b := sl.findPredecessor(key)
	n := b.Next

	count := 0
	for {
		if n == nil {
			break
		}
		c := cpr(key, n.Key)
		if c <= 0 {
			if key&n.Key == key {
				count += len(n.Value)
			}
		}
		b = n
		n = n.Next
	}
	return int32(count)
}

func(n *Node) appendMarker(node *Node){
	n.Next = &Node{
		Key:   0,
		Value: nil,
		Next:  node,
	}
}


func(sl *SkipList) Print(){
	for head := sl.head; head != nil; head = head.down{

		for r := head; r != nil; r = r.right {
			fmt.Print( r.node.Key, "   ")
		}
		fmt.Println()
	}

	for node := sl.head.node; node != nil; node = node.Next {
		fmt.Print(node.Key,"    ")
	}
}

