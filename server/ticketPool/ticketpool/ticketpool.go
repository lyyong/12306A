// @Author: KongLingWen
// @Created at 2021/2/13
// @Modified at 2021/2/13

package ticketpool

import (
	"errors"
	"fmt"
	pb "rpc/ticketPool/proto/ticketPoolRPC"
	"sync"
	"ticketPool/model"
)

type TicketPool struct {
	TrainMap        map[uint32]*Train // key: trainId
	IdToSeatTypeMap map[uint32]string // key: seatTypeId  value: seatTypeName
	SeatTypeToIdMap map[string]uint32
	// 为降低不一致性带来得bug, 将静态资源一起保存
	AllCarriages []*model.CarriageType
	AllSeatInfos map[string]*SeatInfo // 每个车厢对应对应作为类型的SeatInfo, key格式为carriage_id:seatTypeID
}

var (
	TpLock sync.RWMutex
)

type Train struct {
	TrainNum       string
	StopStationMap map[uint32]*StopStation // key: stationId
	CarriageMap    map[string]*Carriages   // key: date  根据日期获得某一天的所有车厢Carriages
}

type StopStation struct {
	Seq         int // 序号：描述该车站是第几个经停站，从 0 开始
	ArriveTime  string
	StartTime   string
	StationName string
}

type Carriages struct {
	CarriageSeatInfo map[uint32]*CarriageSeatInfo // key: seatTypeId  根据 seatTypeId 获取该类型座位的车厢切片
}

type CarriageSeatInfo struct {
	FullValue   uint64
	FullTickets []*FullTicket
	Sl          *SkipList
}

type FullTicket struct {
	Seat              *SeatInfo
	CarriageSeq       string
	CurrentSeatNumber int32  // 从 0 开始，表示当前已分配出去的座位号
	IsAllocate        []bool // 描述数组对应下标编号的座位是否已分配
	Lock              sync.Mutex
}

type SeatInfo struct { // 描述车厢的座位信息，同一种车厢共用一份
	SeatType     string
	MaxSeatCount int32
	Seats        []string           // 票池处理的是整形递增的座位编号，作为下标可以映射为string，如高铁座位的A1 B5...
	ChoseSeatMap map[string][]int32 // key为选择的座位，如'A'，value为这种车厢中'A'座对应的编号
}

func (tp *TicketPool) GetTicket(trainId, startStationId, destStationId uint32, date string, seatCountMap map[uint32]int32) (map[uint32][]string, error) {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	// 根据请求在票池中获取车辆信息，经停站信息，计算 requestValue
	train := tp.TrainMap[trainId]
	if train == nil {
		return nil, errors.New("error train_id")
	}
	startStation := train.StopStationMap[startStationId]
	destStation := train.StopStationMap[destStationId]
	if startStation == nil || destStation == nil || startStation.Seq >= destStation.Seq {
		return nil, errors.New("error station_id")
	}
	requestValue := generateRequestValue(startStation.Seq-1, destStation.Seq-1)
	carriages := train.CarriageMap[date]
	if carriages == nil {
		return nil, errors.New("error date")
	}

	csiNodeMap := make(map[*CarriageSeatInfo]*Node)
	seatsMap := make(map[uint32][]string, len(seatCountMap))

	for seatType, count := range seatCountMap {
		// csi 为描述某种座位余票的结构体
		csi := carriages.CarriageSeatInfo[seatType]
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

func (tp *TicketPool) GetTickets(req *pb.GetTicketRequest) (map[uint32][]string, map[uint32]string, error) {
	// 根据请求在票池中获取车辆信息，经停站信息，计算 requestValue
	train := tp.TrainMap[req.TrainId]
	if train == nil {
		return nil, nil, errors.New("error train_id")
	}
	startStation := train.StopStationMap[req.StartStationId]
	destStation := train.StopStationMap[req.DestStationId]
	if startStation == nil || destStation == nil || startStation.Seq >= destStation.Seq {
		return nil, nil, errors.New("error station_id")
	}
	requestValue := generateRequestValue(startStation.Seq-1, destStation.Seq-1)
	carriages := train.CarriageMap[req.Date]
	if carriages == nil {
		return nil, nil, errors.New("error date")
	}

	choseSeatMap := make(map[uint32]string)
	csiNodeMap := make(map[*CarriageSeatInfo]*Node)
	// 暂存 csi 和已分配的车票（Node），如果后续出票失败，遍历已出车票进行退票
	seatCountMap := make(map[uint32]int32)
	seatsMap := make(map[uint32][]string, len(seatCountMap))

	for i := 0; i < len(req.Passengers); i++ {
		seatTypeId := req.Passengers[i].SeatTypeId
		if req.Passengers[i].ChooseSeat != "" {
			// 处理选座请求
			csi := carriages.CarriageSeatInfo[seatTypeId]
			node := csi.choseSeat(req.Passengers[i].ChooseSeat)
			if node != nil {
				choseSeatMap[req.Passengers[i].PassengerId] = node.Value[0]
				node.Next = csiNodeMap[csi]
				csiNodeMap[csi] = node
				continue
			}
			// node == nil 表示选座失败，以不选座方式进行出票
		}
		seatCountMap[seatTypeId]++
	}

	for seatType, count := range seatCountMap {
		// csi 为描述某种座位余票的结构体
		csi := carriages.CarriageSeatInfo[seatType]
		seatNode, seats := csi.allocateTicket(requestValue, count)
		if seatNode == nil {
			// 有任何一种类型的票出票失败，整个订单都失败
			for c, node := range csiNodeMap {
				for ; node != nil; node = node.Next {
					c.put(node.Key, node.Value)
				}
			}
			return nil, nil, errors.New("there are not enough tickets in the ticket pool")
		}
		// 记录 node 指针、cis指针，如果中途出票失败，则原样退回票池，如果出票成功，则计算余票插入票池
		for n := seatNode; ; n = n.Next {
			// csiNodeMap中可能已有该座位类型的Node(由选座阶段缓存)，因此遍历到该链表尾部进行连接
			if n.Next == nil {
				n.Next = csiNodeMap[csi]
				break
			}
		}
		csiNodeMap[csi] = seatNode
		seatsMap[seatType] = seats
	}

	for csi, node := range csiNodeMap {
		for ; node != nil; node = node.Next {
			remainValue := node.Key ^ requestValue
			csi.put(remainValue, node.Value)
		}
	}
	return seatsMap, choseSeatMap, nil
}

func (tp *TicketPool) SearchTicketCount(trainId, startStationId, destStationId uint32, date string) (map[uint32]int32, error) {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	train := tp.TrainMap[trainId]
	if train == nil {
		return nil, errors.New("error train_id")
	}
	startStation := train.StopStationMap[startStationId]
	destStation := train.StopStationMap[destStationId]
	if startStation == nil || destStation == nil || startStation.Seq >= destStation.Seq {
		return nil, errors.New("error station_id")
	}
	requestValue := generateRequestValue(startStation.Seq-1, destStation.Seq-1)
	carriages := train.CarriageMap[date]
	if carriages == nil {
		return nil, errors.New("error date")
	}

	seatCountMap := make(map[uint32]int32)
	for seatTypeId, csi := range carriages.CarriageSeatInfo {
		seatCountMap[seatTypeId] = csi.getTicketCount(requestValue)
	}

	return seatCountMap, nil
}

func (tp *TicketPool) RefundTickets(trainId, startStationId, destStationId uint32, date string, seatTypeId uint32, seatInfo string) error {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	train := tp.TrainMap[trainId]
	if train == nil {
		return errors.New("error train_id")
	}
	startStation := train.StopStationMap[startStationId]
	destStation := train.StopStationMap[destStationId]
	if startStation == nil || destStation == nil || startStation.Seq >= destStation.Seq {
		return errors.New("error station_id")
	}
	key := generateRequestValue(startStation.Seq-1, destStation.Seq-1)
	carriages := train.CarriageMap[date]
	if carriages == nil {
		return errors.New("error date")
	}
	csi := carriages.CarriageSeatInfo[seatTypeId]
	csi.refund(key, seatInfo)
	return nil
}

func (tp *TicketPool) GetTrainNumber(trainId uint32) string {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	return tp.TrainMap[trainId].TrainNum
}

func (tp *TicketPool) GetSeatTypeNameById(seatTypeId uint32) string {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	return tp.IdToSeatTypeMap[seatTypeId]
}

func (tp *TicketPool) GetIdBySeatTypeName(seatTypeName string) uint32 {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	return tp.SeatTypeToIdMap[seatTypeName]
}

func (csi *CarriageSeatInfo) allocateTicket(requestValue uint64, count int32) (*Node, []string) {
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

func (tp *TicketPool) GetStopStation(trainId, stationId uint32) *StopStation {
	TpLock.RLock()
	defer func() {
		TpLock.RUnlock()
	}()
	return tp.TrainMap[trainId].StopStationMap[stationId]
}

func (csi *CarriageSeatInfo) choseSeat(choice string) *Node {
	for i := 0; i < len(csi.FullTickets); i++ {
		ft := csi.FullTickets[i]
		csm := ft.Seat.ChoseSeatMap
		seatNums := csm[choice]
		for j := 0; j < len(seatNums); j++ {
			seatNum := seatNums[j]
			if !ft.IsAllocate[seatNum] {
				ft.IsAllocate[seatNum] = true
				return &Node{
					Key:   csi.FullValue,
					Value: []string{fmt.Sprintf("%s %s", ft.CarriageSeq, ft.Seat.Seats[seatNum])},
					Next:  nil,
				}
			}
		}
	}
	return nil
}

func (csi *CarriageSeatInfo) splitFullTicket(count int32) *Node {
	seats := make([]string, count)
	num := 0
	for i := 0; i < len(csi.FullTickets); i++ {
		for {
			ft := csi.FullTickets[i]
			maxSeatCount := ft.Seat.MaxSeatCount
			csn := ft.CurrentSeatNumber
			if csn >= maxSeatCount {
				// 当前车厢全票已拆完
				break
			}

			ft.Lock.Lock()
			for count > 0 && ft.CurrentSeatNumber < maxSeatCount {
				n := ft.CurrentSeatNumber
				if !ft.IsAllocate[n] {
					ft.IsAllocate[n] = true
					seats[num] = fmt.Sprintf("%s %s", ft.CarriageSeq, ft.Seat.Seats[n])
					num++
					count--
				}
				ft.CurrentSeatNumber++
			}
			ft.Lock.Unlock()

			//split := maxSeatCount - csn
			//if split >= count {
			//	split = count
			//}
			//if !atomic.CompareAndSwapInt32(&ft.CurrentSeatNumber, csn, csn+split) {
			//	continue
			//}
			//count -= split
			//
			//for j := csn; j < csn+split; j++ {
			//	seats[num] = fmt.Sprintf("%s %s", ft.CarriageSeq, ft.Seat.Seats[j])
			//	num++
			//}
			if count == 0 {
				return &Node{
					Key:   csi.FullValue,
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
		csi.put(csi.FullValue, seats[:len(seats)-int(count)])
	}
	return nil
}

func (csi *CarriageSeatInfo) getTicketCount(requestValue uint64) int32 {
	fullTickets := csi.FullTickets
	var count int32 = 0
	for i := 0; i < len(fullTickets); i++ {
		count += fullTickets[i].Seat.MaxSeatCount - fullTickets[i].CurrentSeatNumber

	}
	// 全票全部被拆分后去票池搜索
	if count == 0 {
		count += csi.search(requestValue)
	}
	return count
}

func (csi *CarriageSeatInfo) getValues(node *Node) []string {
	values := make([]string, 0)
	for ; node != nil; node = node.Next {
		values = append(values, node.Value...)
	}
	return values
}

func (csi *CarriageSeatInfo) refund(key uint64, value string) {
	csi.Sl.Do("Refund", key, value)
}

func (csi *CarriageSeatInfo) allocate(requestValue uint64, count int) *Node {
	respChan := make(chan *Node, 1)
	csi.Sl.Do("Allocate", requestValue, count, respChan)
	return <-respChan
}

func (csi *CarriageSeatInfo) put(key uint64, value []string) {
	csi.Sl.Do("Put", key, value)
}

func (csi *CarriageSeatInfo) search(requestValue uint64) int32 {
	respChan := make(chan int32, 1)
	csi.Sl.Do("Search", requestValue, respChan)
	return <-respChan
}

func generateRequestValue(startStation, destStation int) uint64 {
	var value uint64 = 1
	value <<= destStation - startStation
	value -= 1
	value <<= startStation
	return value
}
