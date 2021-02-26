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
func QueryTrainByDateAndCity(date,startStation, endStation string) []string {

	startCity,_:=RedisDB.HGet("stationCity",startStation).Result()
	endCity,_:=RedisDB.HGet("stationCity",endStation).Result()
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
func QueryTicketNumByDate(date,startStation,endStation string) []*outer.Train {
	//2021-1-23:K4729:1:secondSeat
	trainNos := QueryTrainByDateAndCity(date, startStation, endStation)
	startStationId:=dao.GetStationId(startStation)
	endStationId:=dao.GetStationId(endStation)

	request:=&ticketPoolRPC.GetTicketNumberRequest{}
	request.Date=date
	request.StartStationId=startStationId
	request.DestStationId=endStationId
	for _,trainNo:=range trainNos{
		trainId:=dao.GetTrainId(trainNo)
		request.TrainId=append(request.TrainId,trainId)
	}
	var rpcClient *Client.TPRPCClient
	rpcClient,err=Client.NewClient()
	response,err:=rpcClient.GetTicketNumber(request)
	if err!=nil{
		fmt.Println("rpc getTicketNumber failed, err:",err)
		return nil
	}

	var trains []*outer.Train
	ticketInfos:=response.TrainsTicketInfo
	for i:=0;i<len(ticketInfos);i++{
		ticketInfo:=ticketInfos[i]
		train:=&outer.Train{}
		train.TrainNumber=dao.GetTrainNumber(ticketInfo.TrainId)
		seatInfos:=ticketInfo.SeatInfo
		for _,seatInfo:=range seatInfos{

			switch seatInfo.SeatTypeId {
			case 1:train.BusinessSeat=int(seatInfo.SeatNumber)
			case 2:train.FirstSeat=int(seatInfo.SeatNumber)
			case 3:train.SecondSeat=int(seatInfo.SeatNumber)
			case 4:train.SeniorSoftSleeper=int(seatInfo.SeatNumber)
			case 5:train.SoftSleeper=int(seatInfo.SeatNumber)
			case 6:train.HardSleeper=int(seatInfo.SeatNumber)
			case 7:train.HardSeat=int(seatInfo.SeatNumber)
			default:
			}
		}
		train.TrainType=""
		startStation:=dao.GetStationName(startStationId)
		endStation:=dao.GetStationName(endStationId)
		train.StartStation=startStation
		train.EndStation=endStation

		resMap,_:=RedisDB.HGetAll(trainNos[i]).Result()
		if resMap[startStation]=="1"{
			train.StartStationType="始"
		}else{
			train.StartStationType="终"
		}
		if resMap[endStation]==resMap["stationNum"]{
			train.EndStationType="过"
		}
		leaveTime,_:=RedisDB.HGet(trainNos[i]+":"+resMap[startStation],"leaveTime").Result()
		arriveTime,_:=RedisDB.HGet(trainNos[i]+":"+resMap[endStation],"arriveTime").Result()
		train.StartTime=leaveTime
		train.EndTime=arriveTime

		trains=append(trains,train)
	}
	return trains
}


