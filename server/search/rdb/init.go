/*
* @Author: 余添能
* @Date:   2021/3/3 8:53 下午
 */
package rdb

import (
	"12306A-search/dao"
	"12306A-search/tools/settings"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var RedisDB *redis.Client
var shaBuyTicket string
var err error

func InitRedis() {

	RedisDB = redis.NewClient(&redis.Options{
		Addr:     settings.RedisDB.Host,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	WriteStationAndCityToRedis()
	WriteTrainInfoToRedis()
	WriteTrainPoolToRedis()
}

//用出发城市：目的城市作为key，使用zset类型存储两个城市之间的车次
func WriteTrainPoolToRedis() {
	trainPools := dao.SelectTrainPoolAll()
	if trainPools == nil || len(trainPools) == 0 {
		return
	}
	key := trainPools[0].StartCity + "-" + trainPools[0].EndCity
	exists, err := RedisDB.Exists(key).Result()
	if err != nil {
		fmt.Println("redis.Exists error", err)
		return
	}
	//只要有一趟车次，就默认全部拥有
	if exists > 0 {
		return
	}

	for _, trainPool := range trainPools {
		key := trainPool.StartCity + "-" + trainPool.EndCity

		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", trainPool.StartTime, time.Local)
		//endTime,_:=time.ParseInLocation("2006-01-02 15:04:05", trainPool.EndTime, time.Local)
		//以上车时间作为score,小时+分钟，不用日期，因为车次每天都会有
		hour, minute, _ := startTime.Clock()
		start := hour*60 + minute
		RedisDB.ZAdd(key, redis.Z{Score: float64(start), Member: trainPool.TrainNo})
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

	stopInfos := dao.SelectStopInfoAll()
	key := stopInfos[0].TrainNumber
	exists, err := RedisDB.Exists(key).Result()
	if err != nil {
		fmt.Println("redis.Exists error", err)
		return
	}
	//只要有一趟车次，就默认全部拥有
	if exists > 0 {
		return
	}

	for _, stopInfo := range stopInfos {
		key := stopInfo.TrainNumber
		RedisDB.HSet(key, "stationNum", 0)
	}

	for _, stopInfo := range stopInfos {
		key := stopInfo.TrainNumber

		//车次元数据
		//站名-站序
		RedisDB.HSet(key, stopInfo.StationName, stopInfo.StopSeq)
		//城市-站名
		RedisDB.HSet(key, stopInfo.City, stopInfo.StationName)
		//记录站数
		RedisDB.HIncrBy(key, "stationNum", 1)
		//rdb.RedisDB.HSet(key,stopInfo.City,stopInfo.StationName)

		//保存站点信息:G21-1
		stationKey := stopInfo.TrainNumber + "-" + strconv.Itoa(stopInfo.StopSeq)
		RedisDB.HSet(stationKey, "stopSeq", stopInfo.StopSeq)
		RedisDB.HSet(stationKey, "stationID", stopInfo.StationId)
		RedisDB.HSet(stationKey, "stationName", stopInfo.StationName)
		RedisDB.HSet(stationKey, "cityName", stopInfo.City)
		RedisDB.HSet(stationKey, "arriveTime", stopInfo.ArrivedTime)
		RedisDB.HSet(stationKey, "leaveTime", stopInfo.LeaveTime)
		//保存持续时间
		//rdb.RedisDB.HSet(stationKey, "duration", stopInfo.StayDuration)
		//rdb.RedisDB.HSet(stationKey, "mileage", stopInfo.Mileage)
	}
}

//用来一个hash保存 站点-城市 的映射，
func WriteStationAndCityToRedis() {
	stations := dao.SelectStationAll()
	if stations == nil || len(stations) == 0 {
		return
	}
	key := "stationCity"
	exists, err := RedisDB.Exists(key).Result()
	if err != nil {
		fmt.Println("redis.Exists err:", err)
		return
	}
	if exists > 0 {
		return
	}
	for _, s := range stations {
		RedisDB.HSet(key, s.Name, s.City)
	}
}

//
////上海虹桥
//
////用一个ZAdd保存所有站点
//func WriteStationToRedis() {
//	key := "stations"
//	exists, err := RedisDB.Exists(key).Result()
//	if err!=nil{
//		fmt.Println("redis.Exists err:",err)
//		return
//	}
//	if exists>0{
//		return
//	}
//	stations := dao.SelectStationAll()
//
//	if stations==nil || len(stations)==0{
//		return
//	}
//
//	for _, v := range stations {
//		spell := v.Spell[0]
//		//fmt.Println(v,spell)
//		RedisDB.ZAdd(key, redis.Z{Score: float64(spell), Member: v.Name})
//		//rdb.RedisDB.LPush(key,v.Name)
//	}
//}
//
//// WriteTrainsToRedis 使用Zset保存所有的车次信息
//func WriteTrainsToRedis() {
//
//	trains := dao.SelectTrainAll()
//
//	key := "trains"
//	// res,_:=rdb.RedisDB.HKeys(key).Result()
//	// if res!=nil || len(res)>0{
//	// 	return
//	// }
//	for _, v := range trains {
//		RedisDB.ZAdd(key, redis.Z{Score: float64(v.ID), Member: trains})
//	}
//}
//
