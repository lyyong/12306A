/*
* @Author: 余添能
* @Date:   2021/1/26 3:37 下午
 */
package rdb

import (
	"12306A-search/dao"
	"12306A-search/model/outer"
	"fmt"
	"github.com/go-redis/redis"
	"rpc/ticketPool/Client"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"strconv"
	"strings"
	"time"
)

//按城市查询满足条件的车次，返回车次名
func QueryTrainByDateAndCity(date,startCity, endCity string) []string {

	key:=startCity+"-"+endCity
	now:=time.Now()
	nn:=now.Format("2006-01-02")
	min:=0
	if strings.Compare(date,nn)==0{
		//当天
		h,m,_:=now.Clock()
		min=h*60+m
	}

	res,err := RedisDB.ZRangeByScore(key, redis.ZRangeBy{Min: strconv.Itoa(min), Max: "50000"}).Result()
	if err!=nil{
		fmt.Println("select trains failed, err:",err)
		return nil
	}
	return res
}


//查询某一天某个车次的某种座位等级的车票数，合适的
func QueryTicketNumByDate(date,startCity,endCity string) []*outer.Train {
	//2021-1-23:K4729:1:secondSeat
	//将车站改为城市
	startCity,_=RedisDB.HGet("stationCity",startCity).Result()
	endCity,_=RedisDB.HGet("stationCity",endCity).Result()
	trainNos := QueryTrainByDateAndCity(date, startCity,endCity)

	request:=&ticketPoolRPC.GetTicketNumberRequest{}

	var conditions []*ticketPoolRPC.GetTicketNumberRequest_Condition

	for _,trainNo:=range trainNos{
		trainId:=dao.GetTrainId(trainNo)
		//fmt.Println(trainNo)
		condition:=&ticketPoolRPC.GetTicketNumberRequest_Condition{}
		condition.TrainId=dao.GetTrainId(trainNo)
		startStation,_:=RedisDB.HGet(trainNo,startCity).Result()
		endStation,_:=RedisDB.HGet(trainNo,endCity).Result()
		startStationId:=dao.GetStationId(startStation)
		endStationId:=dao.GetStationId(endStation)
		condition.TrainId=trainId
		condition.StartStationId=startStationId
		condition.DestStationId=endStationId
		conditions=append(conditions,condition)
	}
	request.Date=date
	request.Condition=conditions

	rpcClient, err := Client.NewClient()
	if err!=nil{
		fmt.Println("rpc getTicketNumber failed, err:",err)
		return nil
	}

	response,err:=rpcClient.GetTicketNumber(request)
	ticketInfos:=response.TrainsTicketInfo

	var trains []*outer.Train
	for i:=0;i<len(ticketInfos);i++ {
		train := &outer.Train{}
		ticketInfo := ticketInfos[i]
		trainNo := dao.GetTrainNumber(ticketInfo.TrainId)
		train.TrainNumber = trainNo

		seatInfos := ticketInfo.SeatInfo
		for _, seatInfo := range seatInfos {
			switch seatInfo.SeatTypeId {
			case 1:
				train.BusinessSeat = int(seatInfo.SeatNumber)
			case 2:
				train.FirstSeat = int(seatInfo.SeatNumber)
			case 3:
				train.SecondSeat = int(seatInfo.SeatNumber)
			case 4:
				train.SeniorSoftSleeper = int(seatInfo.SeatNumber)
			case 5:
				train.SoftSleeper = int(seatInfo.SeatNumber)
			case 6:
				train.HardSleeper = int(seatInfo.SeatNumber)
			case 7:
				train.HardSeat = int(seatInfo.SeatNumber)
			default:
			}
		}
		resMap,_:=RedisDB.HGetAll(trainNo).Result()
		startStation:=dao.GetStationName(conditions[i].StartStationId)
		train.StartStation=startStation
		train.StartTime,_=RedisDB.HGet(trainNo+"-"+resMap[train.StartStation],"leaveTime").Result()
		endStation:=dao.GetStationName(conditions[i].DestStationId)
		train.EndStation=endStation
		train.EndTime,_=RedisDB.HGet(trainNo+"-"+resMap[train.EndStation],"arriveTime").Result()

		if resMap[startStation]=="1"{
			train.StartStationType="始"
		}else{
			train.StartStationType="过"
		}
		if resMap[endStation]==resMap["stationNum"]{
			train.EndStationType="终"
		}else{
			train.EndStationType="过"
		}

		trains=append(trains,train)
	}
	return trains
}


