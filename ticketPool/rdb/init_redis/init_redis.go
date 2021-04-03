/*
* @Author: 余添能
* @Date:   2021/1/31 6:37 下午
 */
package init_redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"ticketPool/dao"
	"ticketPool/rdb"
	"time"
)

func InitDataRedis() {

	//初始化trainpool
	fmt.Println("初始化redis：trains")
	WriteTrainPoolToRedis()
	// WriteTrainsToRedis()
	//车站
	//fmt.Println("初始化redis：station")
	//WriteStationToRedis()
	//车站:城市
	fmt.Println("初始化redis：station-city映射")
	WriteStationAndCityToRedis()
	//列车信息
	fmt.Println("开始初始化redis：列车基本信息")
	WriteTrainInfoToRedis()
	//票池
	//fmt.Println("开始初始化redis：票池")
	//WriteTicketPoolToRedis()
	//城市之间的车次

}

//用出发城市：目的城市作为key，使用zset类型存储两个城市之间的车次
func WriteTrainPoolToRedis() {
	trainPools := dao.SelectTrainPoolAll()
	for _, trainPool := range trainPools {
		key := trainPool.StartCity + "-" + trainPool.EndCity

		// res,_:=rdb.RedisDB.Exists(key).Result()
		// if res>0{
		// 	continue
		// }
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", trainPool.StartTime, time.Local)
		//endTime,_:=time.ParseInLocation("2006-01-02 15:04:05", trainPool.EndTime, time.Local)
		//以上车时间作为score,小时+分钟，不用日期，因为车次每天都会有
		hour, minute, _ := startTime.Clock()
		start := hour*60 + minute
		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(start), Member: trainPool.TrainNo})
		//fmt.Println(cmd.Result())

	}
}

//假设所有车次天天会有，所以查车次不用日期
//但可能出现停运情况,特殊考虑

//使用hash存储列车基本信息,key: trainNo 如, G104
//列车信息：列车基本信息、车站信息
//列车基本信息用一个hash保存：key=车次, 元素：trainNo,stationNum, 1--stationNum -> cityName
//每个车站都单独用一个hash保存,key=车次+站序，元素：stationNo,stationName,cityName,arriveTime,departTime,duration,price,mileage
func WriteTrainInfoToRedis() {
	//exists:=make(map[string]bool)
	stopInfos := dao.SelectStopInfoAll()
	for _, stopInfo := range stopInfos {
		key := stopInfo.TrainNumber
		rdb.RedisDB.HSet(key, "stationNum", 0)
	}

	for _, stopInfo := range stopInfos {
		key := stopInfo.TrainNumber

		//车次元数据
		//站名-站序
		rdb.RedisDB.HSet(key, stopInfo.StationName, stopInfo.StopSeq)
		//城市-站名
		rdb.RedisDB.HSet(key, stopInfo.City, stopInfo.StationName)
		//记录站数
		rdb.RedisDB.HIncrBy(key, "stationNum", 1)
		//rdb.RedisDB.HSet(key,stopInfo.City,stopInfo.StationName)

		//保存站点信息:G21-1
		stationKey := stopInfo.TrainNumber + "-" + strconv.Itoa(stopInfo.StopSeq)
		rdb.RedisDB.HSet(stationKey, "stopSeq", stopInfo.StopSeq)
		rdb.RedisDB.HSet(stationKey, "stationID", stopInfo.StationId)
		rdb.RedisDB.HSet(stationKey, "stationName", stopInfo.StationName)
		rdb.RedisDB.HSet(stationKey, "cityName", stopInfo.City)
		//到达时间
		//fmt.Println(stopInfo.ArrivedTime,stopInfo.LeaveTime)

		rdb.RedisDB.HSet(stationKey, "arriveTime", stopInfo.ArrivedTime)
		//出发时间
		rdb.RedisDB.HSet(stationKey, "leaveTime", stopInfo.LeaveTime)
		//保存持续时间
		//rdb.RedisDB.HSet(stationKey, "duration", stopInfo.StayDuration)
		//rdb.RedisDB.HSet(stationKey, "mileage", stopInfo.Mileage)
	}
}

//用来一个hash保存 站点-城市 的映射，
func WriteStationAndCityToRedis() {
	stations := dao.SelectStationAll()
	if stations == nil {
		return
	}
	key := "stationCity"
	// res,_:=rdb.RedisDB.Exists(key).Result()
	// if res>0{
	// 	return
	// }
	for _, s := range stations {
		rdb.RedisDB.HSet(key, s.Name, s.City)
	}
}

//上海虹桥

//用一个List保存所有站点
func WriteStationToRedis() {
	stations := dao.SelectStationAll()

	key := "stations"
	// res,_:=rdb.RedisDB.HKeys(key).Result()
	// if res!=nil || len(res)>0{
	// 	return
	// }
	for _, v := range stations {
		spell := v.Spell[0]
		//fmt.Println(v,spell)
		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(spell), Member: v.Name})
		//rdb.RedisDB.LPush(key,v.Name)
	}
}

// WriteTrainsToRedis 使用List保存所有的车次信息
func WriteTrainsToRedis() {
	trains := dao.SelectTrainAll()

	key := "trains"
	// res,_:=rdb.RedisDB.HKeys(key).Result()
	// if res!=nil || len(res)>0{
	// 	return
	// }
	for _, v := range trains {
		rdb.RedisDB.ZAdd(key, redis.Z{Score: float64(v.ID), Member: trains})
	}
}
