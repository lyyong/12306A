/*
* @Author: 余添能
* @Date:   2021/2/4 10:33 下午
 */
package rdb

import (
	"12306A-search/model/outer"
	"fmt"
	"strconv"
	"strings"
)


//查询某个车次的站序、站点、到达时间、离开时间、时长
func QueryStationByTrainNo(trainNo string) []*outer.Station {
	key:=trainNo
	fmt.Println(key)
	num,_:=RedisDB.HGet(key,"stationNum").Result()
	fmt.Println(num)
	stationNum,_:=strconv.ParseInt(num,10,64)
	var stations []*outer.Station
	for i:=1;i<=int(stationNum);i++{
		stationKey:=trainNo+"-"+strconv.Itoa(i)
		stationMap,err:=RedisDB.HGetAll(stationKey).Result()
		if err!=nil{
			fmt.Println("QueryStationByTrainNo failed, err:",err)
			return nil
		}
		station:=&outer.Station{}
		stationNo,_:=strconv.ParseInt(stationMap["stopSeq"],10,64)
		station.StationNo=int(stationNo)
		station.StationName=stationMap["stationName"]
		station.ArriveTime=stationMap["arriveTime"]
		station.DepartTime=stationMap["leaveTime"]
		//station.Duration=stationMap["duration"]
		arriveTime := strings.Split(stationMap["arriveTime"], ":")
		leaveTime := strings.Split(stationMap["leaveTime"], ":")

		startH, _ := strconv.ParseInt(arriveTime[0], 10, 64)
		startM, _ := strconv.ParseInt(arriveTime[1], 10, 64)
		endH, _ := strconv.ParseInt(leaveTime[0], 10, 64)
		endM, _ := strconv.ParseInt(leaveTime[1], 10, 64)
		minutes := (endH*60 + endM) - (startH*60 + startM)
		station.WaitTime=strconv.Itoa(int(minutes))
		stations=append(stations,station)
	}
	return stations
}

