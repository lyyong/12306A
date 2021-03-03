// @Author: KongLingWen
// @Created at 2021/2/13
// @Modified at 2021/2/13

package ticketpool

import (
	"common/tools/logging"
	"fmt"
	"strconv"
	"strings"
	"ticketPool/model"
	"ticketPool/skiplist"
	"time"
)

var (
	Tp *TicketPool
)

//func init(){
//	// 初始化 db
//	logging.Info("Init DataBase Connections")
//
//	var err error
//	db, err = newMysqlDB()
//	if err != nil {
//		logging.Error("Fail to init DB:", err)
//	}
//	tp = initTicketPool()
//	Close()
//}

const ( // seatTypeID
	BUSINESS_SEAT_ID = iota
	FIRST_SEAT_ID
	SECOND_SEAT_ID
)

var (
	allCarriages []*model.CarriageType
	allSeatInfos map[string]*SeatInfo // 每个车厢对应对应作为类型的SeatInfo, key格式为carriage_id:seatTypeID
)

func InitTicketPool() {
	ticketPool := &TicketPool{
		trainMap:            make(map[uint32]*Train),
		seatTypeMap: 		 initSeatTypeMap(),
	}

	// 初始化车厢座位信息 （所有类型车厢）
	// carriageSeatInfoMap := ticketPool.carriageSeatInfoMap

	genStaticInfo()

	// 暂时只初始化G开头的车辆
	// 得到G开头的车次
	trains := model.GetTrainsByNumberLike("G%")
	//trains:=model.GetTrainsByCondition(map[string]interface{}{"number":"G71"})
	for _, train := range trains {
		// 获得列车类型
		trainType := model.GetTrainTypeByID(train.TrainType)
		// 得到车厢列表
		tcList := strings.Split(trainType.CarriageList, ",")
		// 得到真正的车厢列表
		carriageList := make([]*model.CarriageType, len(tcList))
		for i, tc := range tcList {
			cid, _ := strconv.Atoi(tc)
			carriageList[i] = model.GetCarriageTypesByID(uint(cid))
		}
		// 得到停靠站信息
		stopInfos := model.GetStopInfoByTrainID(train.ID)
		ssm := make(map[uint32]*StopStation)
		for _, stopInfo := range stopInfos {
			ssm[uint32(stopInfo.StationID)] = &StopStation{
				Seq:         stopInfo.StopSeq,
				ArriveTime:  stopInfo.ArrivedTime,
				StartTime:   stopInfo.LeaveTime,
				StationName: stopInfo.StationName,
			}
		}

		// 创建今天开始往后7天的carriageMap
		cm := make(map[string]*Carriages)
		t := time.Now()

		for i := 0; i < 7; i++ {
			date := t.Format("2006-01-02")
			cm[date] = genCarriages(train.ID, date, stopInfos, carriageList)
			t = t.Add(time.Hour * 24)
		}
		ticketPool.trainMap[uint32(train.ID)] = &Train{
			TrainNum:       train.Number,
			stopStationMap: ssm,
			carriageMap:    cm,
		}

	}

	Tp = ticketPool
}

// genStaticInfo 获得车厢信息,生成一些静态数据
func genStaticInfo() {
	allCarriages = model.GetCarriageTypes()
	allSeatInfos = make(map[string]*SeatInfo) // 每个车厢对应对应作为类型的SeatInfo, key格式为carriage_id:seatTypeID
	for _, carriage := range allCarriages {
		if carriage.BusinessSeatNumber > 0 {
			allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, BUSINESS_SEAT_ID)] = &SeatInfo{
				SeatType:     "商务座",
				maxSeatCount: int32(carriage.BusinessSeatNumber),
				seats:        strings.Split(carriage.BusinessSeat, ","),
			}
		}
		if carriage.FirstSeatNumber > 0 {
			allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, FIRST_SEAT_ID)] = &SeatInfo{
				SeatType:     "一等座",
				maxSeatCount: int32(carriage.FirstSeatNumber),
				seats:        strings.Split(carriage.FirstSeat, ","),
			}
		}
		if carriage.SecondSeatNumber > 0 {
			allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, SECOND_SEAT_ID)] = &SeatInfo{
				SeatType:     "二等座",
				maxSeatCount: int32(carriage.SecondSeatNumber),
				seats:        strings.Split(carriage.SecondSeat, ","),
			}
		}
		// TODO 添加其他座位类型
	}
}

func initSeatTypeMap() map[uint32]string {
	seatTypeMap := make(map[uint32]string)
	seatTypeName := []string{"商务座","一等座","二等座","高级软卧","软卧","硬卧","硬座"}
	for i := 0; i < len(seatTypeName); i++ {
		seatTypeMap[uint32(i)] = seatTypeName[i]
	}
	return seatTypeMap
}

func genCarriages(trainId uint, date string, stopInfos []*model.StopInfo, carriageList []*model.CarriageType) *Carriages {
	carriages := &Carriages{carriageSeatInfo: make(map[uint32]*CarriageSeatInfo)}
	fullV := generateFullTicketValue(len(stopInfos))
	// 检查车厢,  将座位类型作为上一级, 暂时只有商务座,一等座和二等座
	business := make([]*FullTicket, 0) // 商务座切片
	first := make([]*FullTicket, 0)    // 一等座
	second := make([]*FullTicket, 0)   // 二等座
	for i, carriage := range carriageList {
		if carriage.BusinessSeatNumber > 0 {
			business = append(business, &FullTicket{
				seat:              allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, BUSINESS_SEAT_ID)],
				carriageSeq:       strconv.Itoa(i),
				currentSeatNumber: 0,
			})
		}
		if carriage.FirstSeatNumber > 0 {
			first = append(first, &FullTicket{
				seat:              allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, FIRST_SEAT_ID)],
				carriageSeq:       strconv.Itoa(i),
				currentSeatNumber: 0,
			})
		}
		if carriage.SecondSeatNumber > 0 {
			second = append(second, &FullTicket{
				seat:              allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, SECOND_SEAT_ID)],
				carriageSeq:       strconv.Itoa(i),
				currentSeatNumber: 0,
			})
		}
		// TODO 添加更多座位类型
	}
	carriages.carriageSeatInfo[BUSINESS_SEAT_ID] = &CarriageSeatInfo{
		fullValue:   fullV,
		fullTickets: business,
		sl:          skiplist.NewSkipList(uint32(trainId), date, BUSINESS_SEAT_ID),
	}
	carriages.carriageSeatInfo[FIRST_SEAT_ID] = &CarriageSeatInfo{
		fullValue:   fullV,
		fullTickets: first,
		sl:          skiplist.NewSkipList(uint32(trainId), date, FIRST_SEAT_ID),
	}
	carriages.carriageSeatInfo[SECOND_SEAT_ID] = &CarriageSeatInfo{
		fullValue:   fullV,
		fullTickets: second,
		sl:          skiplist.NewSkipList(uint32(trainId), date, SECOND_SEAT_ID),
	}
	return carriages
}







// 初始化假数据 - 测试
func InitMockData() {
	// 初始化票池
	Tp = &TicketPool{
		trainMap:            make(map[uint32]*Train),
		seatTypeMap: 		 make(map[uint32]string),
	}
	// 初始化车厢类型
	/*
		   mock数据:
			[ 	carriageTypeId : 0
				carriageType:商务
				maxSeatCount:100
			]
			[ 	carriageTypeId : 1
				carriageType:一等
				maxSeatCount:100
			]
			[ 	carriageTypeId : 2
				carriageType:二等
				maxSeatCount:140
			]
	*/
	seats := make([]string, 100)
	index := 0
	for i := 1; i <= 100/5; i++ {
		s := 'A'
		for j := 0; j < 6; j++ {
			if j == 4 {
				continue
			}
			seat := fmt.Sprintf("%c%d", s+int32(j), i)
			seats[index] = seat
			index++
		}
	}

	businessSeat := &SeatInfo{
		SeatType:     "商务",
		maxSeatCount: 100,
		seats:        seats,
	}
	firstSeat := &SeatInfo{
		SeatType:     "一等座",
		maxSeatCount: 100,
		seats:        seats,
	}

	seatsLevelSecond := make([]string, 140)
	index = 0
	for i := 1; i <= 140/5; i++ {
		s := 'A'
		for j := 0; j < 6; j++ {
			if j == 4 {
				continue
			}
			seat := fmt.Sprintf("%c%d", s+int32(j), i)
			seatsLevelSecond[index] = seat
			index++
		}
	}
	secondSeat := &SeatInfo{
		SeatType:     "二等座",
		maxSeatCount: 140,
		seats:        seatsLevelSecond,
	}

	/* 获取所有列车信息，循环对每一个列车初始化，此处假数据只生成一辆列车
	1.根据列车 id 查询并初始化经停站信息
	2.根据列车 id 查询并初始化车厢信息
	*/

	train := &Train{
		stopStationMap: make(map[uint32]*StopStation),
		carriageMap:    make(map[string]*Carriages),
	}
	Tp.trainMap[0] = train

	t, err := time.Parse("2006-01-02 15:04", "2021-02-16 09:30")
	if err != nil {
		logging.Error("time format error!")
	}
	// 20个站点
	stationNumber := 20
	for i := 0; i < stationNumber; i++ {
		train.stopStationMap[uint32(i)] = &StopStation{
			Seq:        i,
			ArriveTime: t.Format("15:04"),
			StartTime:  t.Add(time.Minute * 10).Format("15:04"),
		}
		t = t.Add(time.Hour)
	}

	// 初始化 2021-02-16 这一天的票
	date, _ := time.Parse("2006-01-02", "2021-02-16")
	carriages := &Carriages{
		carriageSeatInfo: make(map[uint32]*CarriageSeatInfo),
	}
	train.carriageMap[date.Format("2006-01-02")] = carriages

	// 6个商务座车厢
	carriageCount := 6
	index = 1
	business := make([]*FullTicket, carriageCount)
	for i := 0; i < carriageCount; i++ {
		business[i] = &FullTicket{
			seat:              businessSeat,
			carriageSeq:       strconv.Itoa(index),
			currentSeatNumber: 0,
		}
		index++
	}
	// 6个一等座车厢
	first := make([]*FullTicket, carriageCount)
	for i := 0; i < carriageCount; i++ {
		first[i] = &FullTicket{
			seat:              firstSeat,
			carriageSeq:       strconv.Itoa(index),
			currentSeatNumber: 0,
		}
		index++
	}
	// 6个二等座车厢
	second := make([]*FullTicket, carriageCount)
	for i := 0; i < carriageCount; i++ {
		second[i] = &FullTicket{
			seat:              secondSeat,
			carriageSeq:       strconv.Itoa(index),
			currentSeatNumber: 0,
		}
		index++
	}

	fullTicketValue := generateFullTicketValue(stationNumber)
	carriages.carriageSeatInfo[0] = &CarriageSeatInfo{
		fullValue:   fullTicketValue,
		fullTickets: business,
		sl:          skiplist.NewSkipList(0, "2021-02-16", 0),
	}
	carriages.carriageSeatInfo[1] = &CarriageSeatInfo{
		fullValue:   fullTicketValue,
		fullTickets: first,
		sl:          skiplist.NewSkipList(0, "2021-02-16", 1),
	}
	carriages.carriageSeatInfo[2] = &CarriageSeatInfo{
		fullValue:   fullTicketValue,
		fullTickets: second,
		sl:          skiplist.NewSkipList(0, "2021-02-16", 2),
	}
}

// 根据经停站个数，产生fullTicketValue
func generateFullTicketValue(stationNumber int) uint64 {
	var fullTicketValue uint64 = 1
	fullTicketValue <<= stationNumber - 1
	fullTicketValue -= 1
	return fullTicketValue
}

func showTicketPoolInfo() {
	for key, value := range Tp.seatTypeMap {
		fmt.Println("carriageTypeId:[", key, "]; seatInfo:[", value, "]")
	}
	for key, value := range Tp.trainMap {
		fmt.Println("trainId:[", key, "]; train:[", value, "]")
		train := value
		for stationId, station := range train.stopStationMap {
			fmt.Println("stationId:[", stationId, "]; station:[", station, "]")
		}
		for date, carriages := range train.carriageMap {
			fmt.Println("date:[", date, "]; carriages:[", carriages, "]")
			for seatTypeId, csi := range carriages.carriageSeatInfo {
				fmt.Println("seatTypeId:[", seatTypeId, "]; csi:[", csi, "]")
				for _, v := range csi.fullTickets {
					fmt.Println("carriage:", v)
				}
			}
		}
	}
}
