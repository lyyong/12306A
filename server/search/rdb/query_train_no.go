/*
* @Author: 余添能
* @Date:   2021/2/2 2:01 下午
 */
package rdb

import (
	outer2 "12306A/server/search/model/outer"

	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
)

//按城市查询满足条件的车次，返回车次名
func QueryTrainByCity(startCity, endCity string) []string {

	key := startCity + "-" + endCity
	slice := RedisDB.ZRangeByScore(key, redis.ZRangeBy{Min: "0", Max: "1000"})
	res, err := slice.Result()
	if err != nil {
		fmt.Println("redis.zrangebyscores failed,err:", err)
		return nil
	}
	//fmt.Println(res)
	return res
}

//查询满足条件的车次的具体信息
func QueryTrainInfoByTrainNo(trainNo string, startCity, endCity string) *outer2.Train {
	key := trainNo
	mapCmd := RedisDB.HGetAll(key)
	resMap, _ := mapCmd.Result()
	stationNum, _ := strconv.ParseInt(resMap["stationNum"], 10, 64)

	var start, end int
	//寻找对应站点
	for i := 1; i <= (int)(stationNum); i++ {
		if resMap[strconv.Itoa(i)] == startCity {
			start = i
		} else if resMap[strconv.Itoa(i)] == endCity {
			end = i
		} else {

		}
	}
	//寻找上车站和下车站的具体信息
	startKey := trainNo + "-" + strconv.Itoa(start)
	startStation, _ := RedisDB.HGetAll(startKey).Result()

	endKey := trainNo + "-" + strconv.Itoa(end)
	endStation, _ := RedisDB.HGetAll(endKey).Result()
	train := &outer2.Train{}
	train.TrainNo = trainNo
	train.StartStation = startStation["stationName"]
	train.StartTime = startStation["arriveTime"]
	//fmt.Println("sarriveTime",startStation["arriveTime"])
	train.EndStation = endStation["stationName"]
	train.EndTime = endStation["arriveTime"]
	//fmt.Println("earriveTime",endStation["arriveName"])
	//持续时间
	//fmt.Println(startStation["duration"],endStation["duration"])
	startDuration := strings.Split(startStation["duration"], ":")
	endDuration := strings.Split(endStation["duration"], ":")
	startH, _ := strconv.ParseInt(startDuration[0], 10, 64)
	startM, _ := strconv.ParseInt(startDuration[1], 10, 64)
	endH, _ := strconv.ParseInt(endDuration[0], 10, 64)
	endM, _ := strconv.ParseInt(endDuration[1], 10, 64)
	minutes := (endH*60 + endM) - (startH*60 + startM)
	train.Duration = strconv.Itoa(int(minutes/60)) + ":" + strconv.Itoa(int(minutes%60))
	//fmt.Println(train)
	return train
}
