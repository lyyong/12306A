// @Author: KongLingWen
// @Created at 2021/2/13
// @Modified at 2021/2/13

package ticketpool

import (
	"common/tools/logging"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"ticketPool/model"
	"ticketPool/skiplist"
	"time"
)

var (
	Tp   *TicketPool
	Date []*time.Time // 当前可用的日期
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
	logging.Info("Init TicketPool From Serialize File")
	err := InitTicketPoolFromFile()
	if err != nil {
		logging.Error(err)
		logging.Info("Init TicketPool From DB")
		InitTicketPoolFromDB()
	}
	/* 开启票池序列化 */
	Serialize()
	// 开启滚动更新
	rockUpdate(context.TODO())
}

func InitTicketPoolFromFile() error {
	var ticketPool TicketPool
	err := UnSerialize(&ticketPool)
	if err != nil {
		return err
	}
	Tp = &ticketPool
	// 恢复Date
	for _, v := range Tp.TrainMap {
		for k, _ := range v.CarriageMap {
			t, _ := time.Parse("2006-01-02", k)
			Date = append(Date, &t)
		}
		break
	}
	sort.Slice(Date, func(i, j int) bool {
		return Date[i].Before(*Date[j])
	})
	return nil
}

func InitTicketPoolFromDB() {
	ticketPool := &TicketPool{
		TrainMap:        make(map[uint32]*Train),
		IdToSeatTypeMap: initIdToSeatTypeMap(),
		SeatTypeToIdMap: initSeatTypeToIdMap(),
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
			// 添加到时间管理中
			td, _ := time.Parse("2006-01-02", date)
			Date = append(Date, &td)
			cm[date] = genCarriages(train.ID, date, stopInfos, carriageList)
			t = t.Add(time.Hour * 24)
		}
		ticketPool.TrainMap[uint32(train.ID)] = &Train{
			TrainNum:       train.Number,
			StopStationMap: ssm,
			CarriageMap:    cm,
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
				MaxSeatCount: int32(carriage.BusinessSeatNumber),
				Seats:        strings.Split(carriage.BusinessSeat, ","),
			}
		}
		if carriage.FirstSeatNumber > 0 {
			allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, FIRST_SEAT_ID)] = &SeatInfo{
				SeatType:     "一等座",
				MaxSeatCount: int32(carriage.FirstSeatNumber),
				Seats:        strings.Split(carriage.FirstSeat, ","),
			}
		}
		if carriage.SecondSeatNumber > 0 {
			allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, SECOND_SEAT_ID)] = &SeatInfo{
				SeatType:     "二等座",
				MaxSeatCount: int32(carriage.SecondSeatNumber),
				Seats:        strings.Split(carriage.SecondSeat, ","),
			}
		}
		// TODO 添加其他座位类型
	}
}

func initIdToSeatTypeMap() map[uint32]string {
	idToSeatTypeMap := make(map[uint32]string)
	seatTypeNames := []string{"商务座", "一等座", "二等座", "高级软卧", "软卧", "硬卧", "硬座"}
	for i := 0; i < len(seatTypeNames); i++ {
		idToSeatTypeMap[uint32(i)] = seatTypeNames[i]
	}
	return idToSeatTypeMap
}

func initSeatTypeToIdMap() map[string]uint32 {
	seatTypeToIdMap := make(map[string]uint32)
	seatTypeNames := []string{"商务座", "一等座", "二等座", "高级软卧", "软卧", "硬卧", "硬座"}
	for i := 0; i < len(seatTypeNames); i++ {
		seatTypeToIdMap[seatTypeNames[i]] = uint32(i)
	}
	return seatTypeToIdMap
}

func genCarriages(trainId uint, date string, stopInfos []*model.StopInfo, carriageList []*model.CarriageType) *Carriages {
	carriages := &Carriages{CarriageSeatInfo: make(map[uint32]*CarriageSeatInfo)}
	fullV := generateFullTicketValue(len(stopInfos))
	// 检查车厢,  将座位类型作为上一级, 暂时只有商务座,一等座和二等座
	business := make([]*FullTicket, 0) // 商务座切片
	first := make([]*FullTicket, 0)    // 一等座
	second := make([]*FullTicket, 0)   // 二等座
	for i, carriage := range carriageList {
		if carriage.BusinessSeatNumber > 0 {
			business = append(business, &FullTicket{
				Seat:              allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, BUSINESS_SEAT_ID)],
				CarriageSeq:       strconv.Itoa(i),
				CurrentSeatNumber: 0,
			})
		}
		if carriage.FirstSeatNumber > 0 {
			first = append(first, &FullTicket{
				Seat:              allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, FIRST_SEAT_ID)],
				CarriageSeq:       strconv.Itoa(i),
				CurrentSeatNumber: 0,
			})
		}
		if carriage.SecondSeatNumber > 0 {
			second = append(second, &FullTicket{
				Seat:              allSeatInfos[fmt.Sprintf("%d:%d", carriage.ID, SECOND_SEAT_ID)],
				CarriageSeq:       strconv.Itoa(i),
				CurrentSeatNumber: 0,
			})
		}
		// TODO 添加更多座位类型
	}
	carriages.CarriageSeatInfo[BUSINESS_SEAT_ID] = &CarriageSeatInfo{
		FullValue:   fullV,
		FullTickets: business,
		Sl:          skiplist.NewSkipList(uint32(trainId), date, BUSINESS_SEAT_ID),
	}
	carriages.CarriageSeatInfo[FIRST_SEAT_ID] = &CarriageSeatInfo{
		FullValue:   fullV,
		FullTickets: first,
		Sl:          skiplist.NewSkipList(uint32(trainId), date, FIRST_SEAT_ID),
	}
	carriages.CarriageSeatInfo[SECOND_SEAT_ID] = &CarriageSeatInfo{
		FullValue:   fullV,
		FullTickets: second,
		Sl:          skiplist.NewSkipList(uint32(trainId), date, SECOND_SEAT_ID),
	}
	return carriages
}

// 初始化假数据 - 测试
func InitMockData() {
	// 初始化票池
	Tp = &TicketPool{
		TrainMap:        make(map[uint32]*Train),
		IdToSeatTypeMap: make(map[uint32]string),
	}
	// 初始化车厢类型
	/*
		   mock数据:
			[ 	carriageTypeId : 0
				carriageType:商务
				MaxSeatCount:100
			]
			[ 	carriageTypeId : 1
				carriageType:一等
				MaxSeatCount:100
			]
			[ 	carriageTypeId : 2
				carriageType:二等
				MaxSeatCount:140
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
		MaxSeatCount: 100,
		Seats:        seats,
	}
	firstSeat := &SeatInfo{
		SeatType:     "一等座",
		MaxSeatCount: 100,
		Seats:        seats,
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
		MaxSeatCount: 140,
		Seats:        seatsLevelSecond,
	}

	/* 获取所有列车信息，循环对每一个列车初始化，此处假数据只生成一辆列车
	1.根据列车 id 查询并初始化经停站信息
	2.根据列车 id 查询并初始化车厢信息
	*/

	train := &Train{
		StopStationMap: make(map[uint32]*StopStation),
		CarriageMap:    make(map[string]*Carriages),
	}
	Tp.TrainMap[0] = train

	t, err := time.Parse("2006-01-02 15:04", "2021-02-16 09:30")
	if err != nil {
		logging.Error("time format error!")
	}
	// 20个站点
	stationNumber := 20
	for i := 0; i < stationNumber; i++ {
		train.StopStationMap[uint32(i)] = &StopStation{
			Seq:        i,
			ArriveTime: t.Format("15:04"),
			StartTime:  t.Add(time.Minute * 10).Format("15:04"),
		}
		t = t.Add(time.Hour)
	}

	// 初始化 2021-02-16 这一天的票
	date, _ := time.Parse("2006-01-02", "2021-02-16")
	carriages := &Carriages{
		CarriageSeatInfo: make(map[uint32]*CarriageSeatInfo),
	}
	train.CarriageMap[date.Format("2006-01-02")] = carriages

	// 6个商务座车厢
	carriageCount := 6
	index = 1
	business := make([]*FullTicket, carriageCount)
	for i := 0; i < carriageCount; i++ {
		business[i] = &FullTicket{
			Seat:              businessSeat,
			CarriageSeq:       strconv.Itoa(index),
			CurrentSeatNumber: 0,
		}
		index++
	}
	// 6个一等座车厢
	first := make([]*FullTicket, carriageCount)
	for i := 0; i < carriageCount; i++ {
		first[i] = &FullTicket{
			Seat:              firstSeat,
			CarriageSeq:       strconv.Itoa(index),
			CurrentSeatNumber: 0,
		}
		index++
	}
	// 6个二等座车厢
	second := make([]*FullTicket, carriageCount)
	for i := 0; i < carriageCount; i++ {
		second[i] = &FullTicket{
			Seat:              secondSeat,
			CarriageSeq:       strconv.Itoa(index),
			CurrentSeatNumber: 0,
		}
		index++
	}

	fullTicketValue := generateFullTicketValue(stationNumber)
	carriages.CarriageSeatInfo[0] = &CarriageSeatInfo{
		FullValue:   fullTicketValue,
		FullTickets: business,
		Sl:          skiplist.NewSkipList(0, "2021-02-16", 0),
	}
	carriages.CarriageSeatInfo[1] = &CarriageSeatInfo{
		FullValue:   fullTicketValue,
		FullTickets: first,
		Sl:          skiplist.NewSkipList(0, "2021-02-16", 1),
	}
	carriages.CarriageSeatInfo[2] = &CarriageSeatInfo{
		FullValue:   fullTicketValue,
		FullTickets: second,
		Sl:          skiplist.NewSkipList(0, "2021-02-16", 2),
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
	for key, value := range Tp.IdToSeatTypeMap {
		fmt.Println("carriageTypeId:[", key, "]; seatInfo:[", value, "]")
	}
	for key, value := range Tp.TrainMap {
		fmt.Println("trainId:[", key, "]; train:[", value, "]")
		train := value
		for stationId, station := range train.StopStationMap {
			fmt.Println("stationId:[", stationId, "]; station:[", station, "]")
		}
		for date, carriages := range train.CarriageMap {
			fmt.Println("date:[", date, "]; carriages:[", carriages, "]")
			for seatTypeId, csi := range carriages.CarriageSeatInfo {
				fmt.Println("seatTypeId:[", seatTypeId, "]; csi:[", csi, "]")
				for _, v := range csi.FullTickets {
					fmt.Println("carriage:", v)
				}
			}
		}
	}
}

// rockUpdate 滚动更新票池
func rockUpdate(ctx context.Context) {
	// 得到明天凌晨0点的时间点
	updateTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	updateTime = updateTime.Add(24 * time.Hour)
	// 每分钟检查时间是否到0点
	ticker := time.NewTicker(time.Minute)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			if now.Format("2006-01-02") == Date[1].Format("2006-01-02") {
				// 开始初始化
				date := Date[len(Date)-1].Add(24 * time.Hour)
				TpLock.Lock()
				logging.Info("开始更新票池, 开始时间: ", time.Now().Format("2006-01-02; 15:04:05"))
				for trainID, train := range Tp.TrainMap {
					delete(train.CarriageMap, Date[0].Format("2006-01-02"))
					realTrain := model.GetTrainsByNumberLike(train.TrainNum)
					trainType := model.GetTrainTypeByID(realTrain[0].TrainType)
					// 得到车厢列表
					tcList := strings.Split(trainType.CarriageList, ",")
					// 得到真正的车厢列表
					carriageList := make([]*model.CarriageType, len(tcList))
					for i, tc := range tcList {
						cid, _ := strconv.Atoi(tc)
						carriageList[i] = model.GetCarriageTypesByID(uint(cid))
					}
					train.CarriageMap[date.Format("2006-01-02")] = genCarriages(uint(trainID), date.Format("2006-01-02"), model.GetStopInfoByTrainID(uint(trainID)), carriageList)

				}
				Date = append(Date[1:], &date)
				TpLock.Unlock()
				logging.Info("更新票池完成")
			}
		}
	}()
}
