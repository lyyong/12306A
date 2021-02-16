// @Author: KongLingWen
// @Created at 2021/2/13
// @Modified at 2021/2/13

package ticketpool

import (
	"common/tools/logging"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"ticketPool/skiplist"
	"ticketPool/utils/setting"
	"time"
)

var (
	tp *TicketPool
	db *gorm.DB
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

func initTicketPool() *TicketPool {
	ticketPool := &TicketPool{
		trainMap:            make(map[int32]*Train),
		carriageSeatInfoMap: make(map[int32]*SeatInfo),
	}

	// 初始化车厢座位信息 （所有类型车厢）
	//carriageSeatInfoMap := ticketPool.carriageSeatInfoMap


	return ticketPool
}

func newMysqlDB() (*gorm.DB, error) {

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=%s&parseTime=True&loc=Local", setting.DataBase.UserName, setting.DataBase.PassWord, setting.DataBase.DBName, setting.DataBase.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Error("Fail to open db connect:", err)
		return nil, err
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(setting.DataBase.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(setting.DataBase.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	stats, err := json.Marshal(sqlDB.Stats())
	logging.Info("Mysql Connection Pool stats:" + string(stats))

	return db, nil
}

func Close(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func(){
		<-c
		sqlDB,_ := db.DB()
		sqlDB.Close()
		logging.Info("connection pool is closed")
	}()
}


func initMockData(){
	// 连接db

	// 初始化票池
	tp = &TicketPool{
		trainMap:            make(map[int32]*Train),
		carriageSeatInfoMap: make(map[int32]*SeatInfo),
	}
	// 初始化车厢类型
	carriageSeatInfoMap := tp.carriageSeatInfoMap
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
			seat := fmt.Sprintf("%c%d", s + int32(j),i)
			seats[index] = seat
			index++
		}
	}

	carriageSeatInfoMap[0] = &SeatInfo{
		seatType:     "商务",
		maxSeatCount: 100,
		seats:        seats,
	}
	carriageSeatInfoMap[1] = &SeatInfo{
		seatType:     "一等座",
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
			seat := fmt.Sprintf("%c%d", s + int32(j),i)
			seatsLevelSecond[index] = seat
			index++
		}
	}
	carriageSeatInfoMap[2] = &SeatInfo{
		seatType:     "二等座",
		maxSeatCount: 140,
		seats:        seatsLevelSecond,
	}

	/* 获取所有列车信息，循环对每一个列车初始化，此处假数据只生成一辆列车
		1.根据列车 id 查询并初始化经停站信息
		2.根据列车 id 查询并初始化车厢信息
	*/

	train := &Train{
		stopStationMap: make(map[int32]*Station),
		carriageMap:    make(map[string]*Carriages),
	}
	tp.trainMap[0] = train

	t, err := time.Parse("2006-01-02 15:04", "2021-02-16 9:30")
	if err != nil {
		logging.Error("time format error!")
	}
	// 20个站点
	stationNumber := 20
	for i := 0 ; i < stationNumber; i++ {
		train.stopStationMap[int32(i)] = &Station{
			seq:        i,
			arriveTime: t.Format("2006-01-02 15:04"),
			startTime:  t.Add(time.Minute * 10).Format("2006-01-02 15:04"),
		}
		t = t.Add(time.Hour)
	}

	// 初始化 2021-02-16 这一天的票
	date, _ := time.Parse("2006-01-02", "2021-02-16")
	carriages := &Carriages{
		carriageSeatInfo: make(map[int32]*CarriageSeatInfo),
	}
	train.carriageMap[date.Format("2006-01-02")] = carriages

	// 6个商务座车厢
	carriageCount := 6
	index = 1
	business := make([]*FullTicket, carriageCount)
	for i := 0 ; i < carriageCount; i++ {
		business[i] = &FullTicket{
			seat:              tp.carriageSeatInfoMap[0],
			carriageSeq:       strconv.Itoa(index),
			maxSeatCount:      tp.carriageSeatInfoMap[0].maxSeatCount,
			currentSeatNumber: 0,
		}
		index++
	}
	// 6个一等座车厢
	first := make([]*FullTicket, carriageCount)
	for i := 0 ; i < carriageCount; i++ {
		first[i] = &FullTicket{
			seat:              tp.carriageSeatInfoMap[1],
			carriageSeq:       strconv.Itoa(index),
			maxSeatCount:      tp.carriageSeatInfoMap[1].maxSeatCount,
			currentSeatNumber: 0,
		}
		index++
	}
	// 6个二等座车厢
	second := make([]*FullTicket, carriageCount)
	for i := 0 ; i < carriageCount; i++ {
		second[i] = &FullTicket{
			seat:              tp.carriageSeatInfoMap[2],
			carriageSeq:       strconv.Itoa(index),
			maxSeatCount:      tp.carriageSeatInfoMap[2].maxSeatCount,
			currentSeatNumber: 0,
		}
		index++
	}

	fullTicketValue := generateFullTicketValue(stationNumber)
	carriages.carriageSeatInfo[0] = &CarriageSeatInfo{
		fullValue:   fullTicketValue,
		fullTickets: business,
		sl:          skiplist.NewSkipList(),
	}
	carriages.carriageSeatInfo[1] = &CarriageSeatInfo{
		fullValue:   fullTicketValue,
		fullTickets: first,
		sl:          skiplist.NewSkipList(),
	}
	carriages.carriageSeatInfo[2] = &CarriageSeatInfo{
		fullValue:   fullTicketValue,
		fullTickets: second,
		sl:          skiplist.NewSkipList(),
	}

}

// 根据经停站个数，产生fullTicketValue
func generateFullTicketValue(stationNumber int) uint64 {
	var fullTicketValue uint64 = 1
	fullTicketValue <<= stationNumber-1
	fullTicketValue -= 1
	return fullTicketValue
}