// @Author: KongLingWen
// @Created at 2021/2/13
// @Modified at 2021/2/13

package ticketpool

import (
	"errors"
	"fmt"
	"sync/atomic"
	"ticketPool/skiplist"
)

type TicketPool struct {
	trainMap  			map[uint32]*Train 	// key: trainId
	carriageSeatInfoMap map[uint32]*SeatInfo	// key: carriageTypeId
}

type Train struct {
	stopStationMap 	map[uint32]*StopStation // key: stationId
	carriageMap		map[string]*Carriages // key: date  根据日期获得某一天的所有车厢Carriages
}

type StopStation struct{
	Seq        int // 序号：描述该车站是第几个经停站，从 0 开始
	ArriveTime string
	StartTime  string
}

type Carriages struct {
	carriageSeatInfo 	map[uint32]*CarriageSeatInfo // key: seatTypeId  根据 seatTypeId 获取该类型座位的车厢切片
}

type CarriageSeatInfo struct {
	fullValue		uint64
	fullTickets		[]*FullTicket
	sl 				*skiplist.SkipList
}


type FullTicket struct {
	seat				*SeatInfo
	carriageSeq			string
	currentSeatNumber	int32		// 从 0 开始，表示当前已分配出去的座位号
}

type SeatInfo struct {	// 描述车厢的座位信息，同一种车厢共用一份
	seatType 		string
	maxSeatCount 	int32
	seats 			[]string // 票池处理的是整形递增的座位编号，作为下标可以映射为string，如高铁座位的A1 B5...
}



func(tp *TicketPool) GetTicket(trainId, startStationId, destStationId uint32, date string, seatCountMap map[uint32]int32) (map[uint32][]string, error) {
	// 根据请求在票池中获取车辆信息，经停站信息，计算 requestValue
	train := tp.trainMap[trainId]
	startStation := train.stopStationMap[startStationId]
	destStation := train.stopStationMap[destStationId]
	requestValue := generateRequestValue(startStation.Seq, destStation.Seq)
	// 根据日期获取当天的 Carriages
	carriages := train.carriageMap[date]

	csiNodeMap := make(map[*CarriageSeatInfo]*skiplist.Node)
	seatsMap := make(map[uint32][]string,len(seatCountMap))

	for seatType, count := range seatCountMap {
		// csi 为描述某种座位余票的结构体
		csi := carriages.carriageSeatInfo[seatType]
		seatNode, seats := csi.allocateTicket(requestValue, count)
		if seatNode == nil {
			// 有任何一种类型的票出票失败，整个订单都失败
			for c, node := range csiNodeMap {
				for ; node != nil; node = node.Next {
					c.put(node.Key, node.Value)
				}
			}
			return nil, errors.New("there are not enough tickets in the ticket pool")
		}
		// 记录 node 指针、cis指针，如果中途出票失败，则原样退回票池，如果出票成功，则计算余票插入票池
		csiNodeMap[csi] = seatNode
		seatsMap[seatType] = seats
	}

	for csi, node := range csiNodeMap {
		for ; node != nil; node = node.Next {
			remainValue := node.Key ^ requestValue
			csi.put(remainValue, node.Value)
		}
	}
	return seatsMap, nil
}

func(tp *TicketPool) SearchTicketCount(trainId , startStationId, destStationId uint32, date string) map[uint32]int32 {

	train := tp.trainMap[trainId]
	startStation := train.stopStationMap[startStationId]
	destStation := train.stopStationMap[destStationId]
	requestValue := generateRequestValue(startStation.Seq, destStation.Seq)
	carriages := train.carriageMap[date]

	seatCountMap := make(map[uint32]int32)

	for seatTypeId, csi := range carriages.carriageSeatInfo {
		seatCountMap[seatTypeId] = csi.getTicketCount(requestValue)
	}

	return seatCountMap
}

func(tp *TicketPool) GetTrain(trainId uint32) *Train{
	return tp.trainMap[trainId]
}

func(csi *CarriageSeatInfo) allocateTicket(requestValue uint64, count int32)(*skiplist.Node, []string){
	// 优先从票池出票
	node := csi.allocate(requestValue, int(count))

	allocateCount := 0
	for n := node; n != nil; n = n.Next {
		allocateCount += len(n.Value)
	}
	// 判断票池出票数量是否足够
	if allocateCount == int(count) {
		return node, csi.getValues(node)
	}
	// 拆分全票
	splitSeats := csi.splitFullTicket(count - int32(allocateCount))

	// 判断票数是否满足
	if splitSeats != nil {
		splitSeats.Next = node
		node = splitSeats
		if len(node.Value)+allocateCount == int(count) {
			return node, csi.getValues(node)
		}
	}

	// 将已出的票退回票池，全票插入票池（由于这个过程中别的线程有可能也拆分了全票，所以全票不能直接回退，放如票池）
	for ; node != nil; node = node.Next {
		csi.put(node.Key, node.Value)
	}
	return nil, nil
}

func(t *Train) GetStopStation (stationId uint32) *StopStation {
	return t.stopStationMap[stationId]
}


func(csi *CarriageSeatInfo) splitFullTicket(count int32) *skiplist.Node {
	seats := make([]string, count)
	num := 0
	for i := 0; i < len(csi.fullTickets); i++ {
		for {
			ft := csi.fullTickets[i]
			maxSeatCount := ft.seat.maxSeatCount
			csn := ft.currentSeatNumber
			if csn >= maxSeatCount {
				// 当前车厢全票已拆完
				break
			}
			split := maxSeatCount - csn
			if split >= count {
				split = count
			}
			if !atomic.CompareAndSwapInt32(&ft.currentSeatNumber, csn, csn+split) {
				continue
			}
			count -= split

			for j := csn; j < csn+split; j++ {
				seats[num] = fmt.Sprintf("%s %s",ft.carriageSeq, ft.seat.seats[j])
				num++
			}
			if count == 0 {
				return &skiplist.Node{
					Key:   csi.fullValue,
					Value: seats,
					Next:  nil,
				}
			}
			// 当前车厢剩余全票少于 count，跳出当前循环继续分配下一个车厢
			break
		}
	}
	// 所有车厢全票之和不足count，返回nil，分配出的票插入票池（有可能别的线程也进行了分配，所以不能再减回去）
	if len(seats) > int(count) {
		csi.put(csi.fullValue, seats[:len(seats)-int(count)])
	}
	return nil
}

func(csi *CarriageSeatInfo) getTicketCount(requestValue uint64) int32 {
	fullTickets := csi.fullTickets
	var count int32 = 0
	for i := 0; i < len(fullTickets); i++{
		count += fullTickets[i].seat.maxSeatCount - fullTickets[i].currentSeatNumber

	}
	// 全票全部被拆分后去票池搜索
	if count == 0 {
		count += csi.search(requestValue)
	}
	return count
}

func (csi *CarriageSeatInfo) getValues(node *skiplist.Node) []string{
	values := make([]string, 0)
	for ; node != nil; node = node.Next {
		values = append(values, node.Value...)
	}
	return values
}

func (csi *CarriageSeatInfo) allocate(requestValue uint64, count int) *skiplist.Node{
	respChan := make(chan *skiplist.Node, 1)
	csi.sl.Do("Allocate", requestValue, count, respChan)
	return <- respChan
}

func (csi *CarriageSeatInfo) put(key uint64, value []string){
	csi.sl.Do("Put", key, value)
}

func (csi *CarriageSeatInfo) search(requestValue uint64) int32 {
	respChan := make(chan int32, 1)
	csi.sl.Do("Search", requestValue, respChan)
	return <- respChan
}

func generateRequestValue(startStation, destStation int) uint64{
	var value uint64 = 1
	value <<= destStation - startStation
	value -= 1
	value <<= startStation
	return value
}


