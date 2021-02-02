/*
* @Author: 余添能
* @Date:   2021/1/31 6:37 下午
 */
package rdb

import (
	"12306A/ticketPool/init_data"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var RedisDB *redis.Client

func init() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

//假设所有车次天天会有，所以查车次不用日期
//但可能出现停运情况,特殊考虑

//用出发城市：目的城市作为key，使用list类型存储两个城市之间的车次
func WriteTrainPoolToRedis() {
	trainPools := init_data.ReadTrainPoolAll()
	for _, trainPool := range trainPools {
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", trainPool.StartTime, time.Local)
		//endTime,_:=time.ParseInLocation("2006-01-02 15:04:05", trainPool.EndTime, time.Local)
		//以上车时间作为score,小时+分钟，不用日期，因为车次每天都会有
		hour, minute, _ := startTime.Clock()
		start := hour*60 + minute
		//fmt.Println(start)
		key := trainPool.StartCity + "-" + trainPool.EndCity
		RedisDB.ZAdd(key, redis.Z{Score: float64(start), Member: trainPool.TrainNo})
		//fmt.Println(cmd.Result())

	}
	zslice := RedisDB.ZRangeByScoreWithScores("景德镇九江", redis.ZRangeBy{Min: "0", Max: "50000"})
	fmt.Println(zslice.Val())
}

//使用hash存储列车基本信息,key: trainNo 如, G104
func WriteTicketPoolToRedis() {
	trains := init_data.ReadTotalTrainNo()
	for _, train := range trains {
		//initialTime:=train.InitialTime
		trainNo := train.TrainNo
		key := trainNo
		//写入车次元数据
		stationNum := train.StationNum
		RedisDB.HSet(key, "stationNum", stationNum)
		//RedisDB.HSet(key,"initialTime",train.InitialTime)
		//RedisDB.HSet(key,"terminalTime",train.TerminalTime)
		stations := train.Stations
		for i := 0; i < len(stations); i++ {
			//保存各站所在城市
			RedisDB.HSet(key, strconv.Itoa(i+1), stations[i].CityName)

			//保存站点信息
			stationKey := trainNo + "-" + strconv.Itoa(i+1)
			RedisDB.HSet(stationKey, "stationNo", stations[i].StationNo)
			RedisDB.HSet(stationKey, "stationName", stations[i].StationName)
			RedisDB.HSet(stationKey, "cityName", stations[i].CityName)
			//到达时间
			hour1, minute1, _ := stations[i].ArriveTime.Clock()
			RedisDB.HSet(stationKey, "arriveTime", strconv.Itoa(hour1)+":"+strconv.Itoa(minute1))
			//出发时间
			hour2, minute2, _ := stations[i].DepartTime.Clock()
			RedisDB.HSet(stationKey, "departTime", strconv.Itoa(hour2)+":"+strconv.Itoa(minute2))
			//保存持续时间
			duration := stations[i].Duration.Sub(stations[0].Duration)
			hour := (int)(duration.Hours())
			minute := (int)(duration.Minutes()) % 60
			RedisDB.HSet(stationKey, "duration", strconv.Itoa(hour)+":"+strconv.Itoa(minute))
			RedisDB.HSet(stationKey, "mileage", stations[i].Mileage)
			RedisDB.HSet(stationKey, "price", stations[i].Price)
		}
	}
}

//用来一个hash保存 站点-城市 的映射，
func WriteStationAndCityToRedis() {
	sqlStr := "select city,station_name from station_province_city;"
	rows, err := init_data.Db.Query(sqlStr)
	if err != nil {
		fmt.Println("query station_province_city failed, err:", err)
		return
	}
	key := "stationCity"
	for rows.Next() {
		var city, stationName string
		rows.Scan(&city, &stationName)
		RedisDB.HSet(key, stationName, city)
	}
}
