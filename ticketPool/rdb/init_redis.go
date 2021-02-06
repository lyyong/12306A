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
	RedisDB=redis.NewClient(&redis.Options{
		Addr: "0.0.0.0:6379",
	})
	//RedisDB = redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs:  []string{"192.168.10.11:7001","192.168.10.11:7002", "192.168.10.11:7003"},
	//})
	//连接redis集群

}
func InitDataRedis()  {
	//车站
	WriteStationToRedis()
	//车站:城市
	WriteStationAndCityToRedis()
	//列车信息
	WriteTrainInfoToRedis()
	//票池
	WriteTicketPoolToRedis()
	//城市之间的车次
	WriteTrainPoolToRedis()
}
//假设所有车次天天会有，所以查车次不用日期
//但可能出现停运情况,特殊考虑

//用出发城市：目的城市作为key，使用zset类型存储两个城市之间的车次
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
//列车信息：列车基本信息、车站信息
//列车基本信息用一个hash保存：key=车次, 元素：trainNo,stationNum, 1--stationNum -> cityName
//每个车站都单独用一个hash保存,key=车次+站序，元素：stationNo,stationName,cityName,arriveTime,departTime,duration,price,mileage
func WriteTrainInfoToRedis() {
	trains := init_data.ReadTotalTrainNo()
	for _, train := range trains {
		//initialTime:=train.InitialTime
		trainNo := train.TrainNo
		key := trainNo
		//写入车次元数据
		stationNum := train.StationNum
		RedisDB.HSet(key, "stationNum", stationNum)
		//RedisDB.HSet(key, "initialTime", stationNum)
		//RedisDB.HSet(key,"initialTime",train.InitialTime)
		//RedisDB.HSet(key,"terminalTime",train.TerminalTime)
		stations := train.Stations
		for i := 0; i < len(stations); i++ {
			//保存各站所在城市
			RedisDB.HSet(key, strconv.Itoa(i+1), stations[i].CityName)
			//保存各站对应站序
			RedisDB.HSet(key,stations[i].StationName,strconv.Itoa(i+1))
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
//用一个List保存所有站点
func WriteStationToRedis()  {
	sqlStr := "select station_name from station_province_city;"
	rows, err := init_data.Db.Query(sqlStr)
	if err != nil {
		fmt.Println("query station_province_city failed, err:", err)
		return
	}
	key := "stations"
	for rows.Next() {
		var stationName string
		rows.Scan( &stationName)
		RedisDB.LPush(key,stationName)
	}
}

//初始化每趟车次的票池
//second,first,business
//zset保存票
//key=日期:车次::
func WriteTicketPoolToRedis()  {
	trainNos:=init_data.ReadTotalTrainNo()

	for i:=0;i<len(trainNos);i++{
		train:=trainNos[i]
		//fmt.Println(train)
		stations:=train.Stations
		year:=train.InitialTime.Year()
		month:=train.InitialTime.Month()
		day:=train.InitialTime.Day()
		date:=strconv.Itoa(year)+"-"+strconv.Itoa(int(month))+"-"+strconv.Itoa(day)
		for j:=0;j<len(stations);j++{
			stationKey1:=date+":"+train.TrainNo+":"+strconv.Itoa(j+1)+":"+"firstSeat"
			stationKey2:=date+":"+train.TrainNo+":"+strconv.Itoa(j+1)+":"+"secondSeat"
			stationKey3:=date+":"+train.TrainNo+":"+strconv.Itoa(j+1)+":"+"businessSeat"
			//score=下车站序，member=车厢号:座位号
			//需要改
			for k:=1;k<=20;k++{
				RedisDB.ZAdd(stationKey1,redis.Z{Score: float64(len(stations)),Member: strconv.Itoa(k)+":"+strconv.Itoa(i+1)})
				RedisDB.ZAdd(stationKey2,redis.Z{Score: float64(len(stations)),Member: strconv.Itoa(k)+":"+strconv.Itoa(i+1)})
				RedisDB.ZAdd(stationKey3,redis.Z{Score: float64(len(stations)),Member: strconv.Itoa(k)+":"+strconv.Itoa(i+1)})
				//RedisDB.ZRem(stationKey1,strconv.Itoa(k)+":"+strconv.Itoa(i+1))
				//RedisDB.ZRem(stationKey2,strconv.Itoa(k)+":"+strconv.Itoa(i+1))
				//RedisDB.ZRem(stationKey3,strconv.Itoa(k)+":"+strconv.Itoa(i+1))
			}
		}
	}
}

