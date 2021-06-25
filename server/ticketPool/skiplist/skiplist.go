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
	Head        *Index
	MaxLevel    int
	RequestChan chan *Request
	TrainId     uint32
	Date        string
	SeatTypeId  uint32
}

type Node struct {
	Key   uint64
	Value []string
	Next  *Node
}

type Index struct {
	Node  *Node
	Down  *Index
	Right *Index
}

type Request struct {
	Option string
	Args   []interface{}
}

func NewSkipList(trainId uint32, date string, seatTypeId uint32) *SkipList {
	head := &Index{
		Node:  &Node{},
		Down:  nil,
		Right: nil,
	}
	sl := &SkipList{
		Head:        head,
		MaxLevel:    1,
		RequestChan: make(chan *Request, 100),
		TrainId:     trainId,
		Date:        date,
		SeatTypeId:  seatTypeId,
	}
	sl.DealWithRequest()
	return sl
}

func (sl *SkipList) Do(cmd string, args ...interface{}) {
	sl.RequestChan <- &Request{
		Option: cmd,
		Args:   args,
	}
}

func (sl *SkipList) DealWithRequest() {
	go func() {
		for {
			req := <-sl.RequestChan
			switch req.Option {
			case "Put":
				key := req.Args[0].(uint64)
				value := req.Args[1].([]string)
				if key == 0 {
					continue
				}
				sl.Put(key, value)
				persistence.Do(&persistence.PersistentRequest{
					Option:     "INSERT",
					TrainId:    sl.TrainId,
					Date:       sl.Date,
					SeatTypeId: sl.SeatTypeId,
					Key:        key,
					Value:      value,
				})
			case "Allocate":
				key := req.Args[0].(uint64)
				count := req.Args[1].(int)
				respChan := req.Args[2].(chan *Node)
				node := sl.Allocate(key, count)
				respChan <- node
				for ; node != nil; node = node.Next {
					persistence.Do(&persistence.PersistentRequest{
						Option:     "DELETE",
						TrainId:    sl.TrainId,
						Date:       sl.Date,
						SeatTypeId: sl.SeatTypeId,
						Key:        node.Key,
						Value:      node.Value,
					})
				}
			case "Search":
				key := req.Args[0].(uint64)
				respChan := req.Args[1].(chan int32)
				respChan <- sl.Search(key)
			case "Refund":
				key := req.Args[0].(uint64)
				value := req.Args[1].(string)
				sl.refund(key, value)
			}
		}
	}()
}

func (sl *SkipList) Get(key uint64) []string {
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

func (sl *SkipList) Put(key uint64, value []string) {
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
	if rnd&0x80000001 == 0 {
		level := 1
		for rnd >>= 1; (rnd & 1) != 0; rnd >>= 1 {
			level++
		}
		var idx *Index

		// 随机level <= 当前 level 时，生成随机 level 个索引节点
		if level <= sl.MaxLevel {
			for i := 0; i < level; i++ {
				idx = &Index{
					Node:  node,
					Down:  idx,
					Right: nil,
				}
			}
		} else {
			// 随机 level 大于当前 level 时
			// 增加索引层数
			level = sl.MaxLevel + 1
			// 生成当 level 个索引节点
			for i := 0; i < level; i++ {
				idx = &Index{
					Node:  node,
					Down:  idx,
					Right: nil,
				}
			}
			// 增加原头索引高度
			h := sl.Head

			newH := &Index{
				Node:  h.Node,
				Down:  h,
				Right: nil,
			}
			sl.Head = newH
			sl.MaxLevel++
		}

		// 寻找索引插入点
		// 从头结点开始（头结点 level 为最高 level），每层向右寻找插入位置，找到后插入，然后进入下一层
		q := sl.Head
		r := q.Right
		hLevel := sl.MaxLevel
		for {

			if r != nil {
				n := r.Node
				c := cpr(key, n.Key)
				if c > 0 {
					q = r
					r = r.Right
					continue
				}
			}
			if hLevel == level {
				idx.Right = r
				q.Right = idx
				idx = idx.Down
				if level--; level < 1 {
					break
				}
			}
			hLevel--

			q = q.Down
			r = q.Right
		}
	}

}

func (sl *SkipList) Remove(key uint64) []string {
	b := sl.findPredecessor(key)
	n := b.Next

	for {
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
		if sl.Head.Right == nil {
			sl.tryReduceLevel()
		}
		return v
	}
}

func (sl *SkipList) refund(key uint64, value string) {
	h := sl.Head

	for {
		if h.Down == nil {
			break
		}
		h = h.Down
	}
	n := h.Node
	for {
		n = n.Next

		if n == nil {
			sl.Put(key, []string{value})
			return
		}

		values := n.Value
		for i := 0; i < len(values); i++ {
			if values[i] == value {
				newKey := n.Key | key
				if len(values) == 1 {
					sl.Remove(n.Key)
				} else {
					n.Value = append(values[:i], values[i+1:]...)
				}
				sl.Put(newKey, []string{value})
				return
			}
		}
	}
}

func (sl *SkipList) tryReduceLevel() {
	if sl.MaxLevel <= 3 {
		return
	}
	d := sl.Head.Down

	if d == nil && d.Down == nil {
		sl.Head = d
	}
}

// return Node  -->  Node.Key <= Key < Node.Next.Key
func (sl *SkipList) findPredecessor(key uint64) *Node {
	var q, r, d *Index
	q = sl.Head
	r = q.Right
	for {
		if r != nil {
			n := r.Node
			k := n.Key
			if n.Value == nil { // clean index
				q.Right = r.Right
				r = q.Right
				continue
			}
			if cpr(key, k) > 0 { // 与 q 的右节点的 Key 比较，Key > k 则向右寻找
				q = r
				r = r.Right
				continue
			}
		}
		// r == nil 或 Key <= k    r = q.Right n == r.Node
		if d = q.Down; d == nil { // 进入下一层
			return q.Node
		}
		q = d
		r = d.Right
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

func (sl *SkipList) Allocate(key uint64, count int) *Node {
	b := sl.findPredecessor(key)
	n := b.Next

	var node *Node

	for {
		if n == nil {
			return node
		}
		c := cpr(key, n.Key)
		if c <= 0 {
			if key&n.Key == key {
				v := n.Value
				lenV := len(n.Value)
				if lenV > count {
					// 当前节点的座位数量满足需求，分配并返回
					retV := make([]string, count)
					copy(retV, v[lenV-count:])
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
					if sl.Head.Right == nil {
						sl.tryReduceLevel()
					}

					if count == 0 {
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

func (sl *SkipList) Search(key uint64) int32 {
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

func (n *Node) appendMarker(node *Node) {
	n.Next = &Node{
		Key:   0,
		Value: nil,
		Next:  node,
	}
}

func (sl *SkipList) Print() {
	for head := sl.Head; head != nil; head = head.Down {

		for r := head; r != nil; r = r.Right {
			fmt.Print(r.Node.Key, "   ")
		}
		fmt.Println()
	}

	for node := sl.Head.Node; node != nil; node = node.Next {
		fmt.Print(node.Key, "    ")
	}
}
