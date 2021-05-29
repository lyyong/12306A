/*
* @Author: 余添能
* @Date:   2021/1/26 3:37 下午
 */
package rdb

import (
	"12306A-search/dao"
	"12306A-search/model/outer"
	"12306A-search/tools/settings"
	"fmt"
	"github.com/go-redis/redis"
	"rpc/ticketPool/Client"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"strconv"
	"strings"
	"time"
)

//按城市查询满足条件的车次，返回车次名
func QueryTrainByDateAndCity(date, startCity, endCity string) []string {

	key := startCity + "-" + endCity
	now := time.Now()
	nn := now.Format("2006-01-02")
	min := 0
	if strings.Compare(date, nn) == 0 {
		//当天
		h, m, _ := now.Clock()
		//fmt.Println(now.Clock())
		min = h*60 + m
	}
	//fmt.Println(min)
	res, err := RedisDB.ZRangeByScore(key, redis.ZRangeBy{Min: strconv.Itoa(min), Max: "50000"}).Result()
	if err != nil {
		fmt.Println("select trains failed, err:", err)
		return nil
	}
	//for _,t:=range res{
	//	fmt.Println(t)
	//}
	return res
}

//查询某一天某个车次的某种座位等级的车票数，合适的
func QueryTicketNumByDate(date, startCity, endCity string) []*outer.Train {
	//2021-1-23:K4729:1:secondSeat
	//将车站改为城市
	startCity, _ = RedisDB.HGet("stationCity", startCity).Result()
	endCity, _ = RedisDB.HGet("stationCity", endCity).Result()
	trainNos := QueryTrainByDateAndCity(date, startCity, endCity)
	//fmt.Println(trainNos)
	var trains []*outer.Train

	request := &ticketPoolRPC.GetTicketNumberRequest{}

	var conditions []*ticketPoolRPC.GetTicketNumberRequest_Condition

	for _, trainNo := range trainNos {
		trainId := dao.GetTrainId(trainNo)
		//fmt.Println(trainNo)
		condition := &ticketPoolRPC.GetTicketNumberRequest_Condition{}
		condition.TrainId = dao.GetTrainId(trainNo)
		startStation, _ := RedisDB.HGet(trainNo, startCity).Result()
		endStation, _ := RedisDB.HGet(trainNo, endCity).Result()
		startStationId := dao.GetStationId(startStation)
		endStationId := dao.GetStationId(endStation)
		condition.TrainId = trainId
		condition.StartStationId = startStationId
		condition.DestStationId = endStationId
		conditions = append(conditions, condition)
	}
	request.Date = date
	request.Condition = conditions

	var response = &ticketPoolRPC.GetTicketNumberResponse{}
	var trainTicketInfos []*ticketPoolRPC.TrainTicketInfo

	//先查询缓存
	cacheKey := startCity + "-" + endCity + "-cache"
	exists, err := RedisDB.Exists(cacheKey).Result()
	if err != nil {
		return nil
	}
	if exists == 1 {
		fmt.Println("有缓存...")
		date := request.Date
		for _, condition := range request.Condition {
			key := date + "-" + strconv.Itoa(int(condition.TrainId)) +
				"-" + strconv.Itoa(int(condition.StartStationId)) + "-" + strconv.Itoa(int(condition.DestStationId))
			result, err := RedisDB.HGetAll(key).Result()
			if err != nil {
				fmt.Println("err:", err)
				return trains
			}

			var seatInfos []*ticketPoolRPC.SeatInfo
			//座位等级=0，1，2...6
			for k, v := range result {
				seatType, err := strconv.ParseInt(k, 10, 32)
				seatNum, err := strconv.ParseInt(v, 10, 32)
				if err != nil {
					fmt.Println("parseInt error:", err)
					return nil
				}
				seatInfo := &ticketPoolRPC.SeatInfo{}
				seatInfo.SeatTypeId = uint32(seatType)
				seatInfo.SeatNumber = int32(seatNum)
				seatInfos = append(seatInfos, seatInfo)

			}
			//存储一趟车次的余票情况
			trainTicketInfo := &ticketPoolRPC.TrainTicketInfo{}
			trainTicketInfo.TrainId = condition.TrainId
			trainTicketInfo.SeatInfo = seatInfos

			trainTicketInfos = append(trainTicketInfos, trainTicketInfo)
		}
		response.TrainsTicketInfo = trainTicketInfos
	} else {
		fmt.Println("没有缓存")
		//没有缓存
		rpcClient, err := Client.NewClientWithTarget(settings.Target.Addr)
		if err != nil {
			fmt.Println("rpc getTicketNumber failed, err:", err)
			return nil
		}

		response, err = rpcClient.GetTicketNumber(request)
		if response == nil {
			return nil
		}

		//更新或者保存缓存
		RedisDB.Set(cacheKey, "abc", time.Minute*1)

		for k, trainTicketInfo := range response.TrainsTicketInfo {

			trainId := trainTicketInfo.TrainId
			startStationId := request.Condition[k].StartStationId
			endStationId := request.Condition[k].DestStationId

			//2021-05-29-G21-1-5
			cacheKey := date + "-" + strconv.Itoa(int(trainId)) +
				"-" + strconv.Itoa(int(startStationId)) + "-" + strconv.Itoa(int(endStationId))
			seatInfos := trainTicketInfo.SeatInfo

			for _, seatInfo := range seatInfos {
				//存储每一等级座位数量
				RedisDB.HSet(cacheKey, strconv.Itoa(int(seatInfo.SeatTypeId)), seatInfo.SeatNumber)
			}
			//过期时间两分钟
			RedisDB.Expire(cacheKey, time.Minute*1)
		}
	}

	ticketInfos := response.TrainsTicketInfo
	//fmt.Println(ticketInfos)
	if ticketInfos == nil || len(ticketInfos) == 0 {
		return nil
	}

	for i := 0; i < len(ticketInfos); i++ {
		train := &outer.Train{}
		ticketInfo := ticketInfos[i]
		trainNo := dao.GetTrainNumber(ticketInfo.TrainId)
		train.TrainNumber = trainNo
		train.TrainID = uint64(ticketInfo.TrainId)
		train.TrainType = "G"

		seatInfos := ticketInfo.SeatInfo
		for _, seatInfo := range seatInfos {
			switch seatInfo.SeatTypeId {
			case 0:
				train.BusinessSeat = int(seatInfo.SeatNumber)
				train.BusinessSeatPrice = 500
			case 1:
				train.FirstSeat = int(seatInfo.SeatNumber)
				train.FirstSeatPrice = 500
			case 2:
				train.SecondSeat = int(seatInfo.SeatNumber)
				train.SecondSeatPrice = 500
			case 3:
				train.SeniorSoftSleeper = int(seatInfo.SeatNumber)
				train.SeniorSoftBerthPrice = 500
			case 4:
				train.SoftSleeper = int(seatInfo.SeatNumber)
				train.SoftBerthPrice = 500
			case 5:
				train.HardSleeper = int(seatInfo.SeatNumber)
				train.HardBerthPrice = 500
			case 6:
				train.HardSeat = int(seatInfo.SeatNumber)
				train.HardSeatPrice = 500
			default:
			}
		}
		resMap, _ := RedisDB.HGetAll(trainNo).Result()
		leaveStation := dao.GetStationName(conditions[i].StartStationId)
		train.LeaveStation = leaveStation
		train.LeaveStationNo = uint64(conditions[i].StartStationId)
		train.LeaveTime, _ = RedisDB.HGet(trainNo+"-"+resMap[train.LeaveStation], "leaveTime").Result()
		arrivalStation := dao.GetStationName(conditions[i].DestStationId)
		train.ArrivalStation = arrivalStation
		train.ArrivalStationNo = uint64(conditions[i].DestStationId)
		train.ArrivalTime, _ = RedisDB.HGet(trainNo+"-"+resMap[train.ArrivalStation], "arriveTime").Result()

		if resMap[leaveStation] == "1" {
			train.LeaveStationType = "始"
		} else {
			train.LeaveStationType = "过"
		}
		if resMap[arrivalStation] == resMap["stationNum"] {
			train.ArrivalStationType = "终"
		} else {
			train.ArrivalStationType = "过"
		}

		// 获取始发站和终点站
		train.StartStation, _ = RedisDB.HGet(trainNo+"-"+"1", "stationName").Result()
		train.StartStationID, _ = RedisDB.HGet(trainNo+"-"+"1", "stationID").Result()
		train.EndStation, _ = RedisDB.HGet(trainNo+"-"+resMap["stationNum"], "stationName").Result()
		train.EndStationID, _ = RedisDB.HGet(trainNo+"-"+resMap["stationNum"], "stationID").Result()

		trains = append(trains, train)
	}
	return trains
}

// QueryTicketNumByDateWithTrainNumber 查询某车次,具体两站之间的余票数
func QueryTicketNumByDateWithTrainNumber(TrainId, ssID, esID uint32, date string) *outer.Train {
	request := &ticketPoolRPC.GetTicketNumberRequest{}

	conditions := []*ticketPoolRPC.GetTicketNumberRequest_Condition{
		{TrainId: TrainId, StartStationId: ssID, DestStationId: esID},
	}
	request.Date = date
	request.Condition = conditions

	var response = &ticketPoolRPC.GetTicketNumberResponse{}
	var trainTicketInfos []*ticketPoolRPC.TrainTicketInfo

	//先查缓存
	//cacheKey=2021-05-29-G21-1-5
	cacheKey := date + "-" + strconv.Itoa(int(TrainId)) + "-" + strconv.Itoa(int(ssID)) + "-" + strconv.Itoa(int(esID))
	resMap, err := RedisDB.HGetAll(cacheKey).Result()
	if err != nil {
		fmt.Println("Redis.Exists err:", err)
		return nil
	}
	//有缓存
	if resMap != nil && len(resMap) > 0 {
		fmt.Println("有缓存")
		var seatInfos []*ticketPoolRPC.SeatInfo

		for k, v := range resMap {
			seatType, err := strconv.ParseInt(k, 10, 32)
			seatNum, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				fmt.Println("parseInt error:", err)
				return nil
			}
			seatInfo := &ticketPoolRPC.SeatInfo{}
			seatInfo.SeatTypeId = uint32(seatType)
			seatInfo.SeatNumber = int32(seatNum)
			seatInfos = append(seatInfos, seatInfo)
		}
		//存储一趟车次的余票情况
		trainTicketInfo := &ticketPoolRPC.TrainTicketInfo{}
		trainTicketInfo.TrainId = TrainId
		trainTicketInfo.SeatInfo = seatInfos
		trainTicketInfos = append(trainTicketInfos, trainTicketInfo)
		response.TrainsTicketInfo = trainTicketInfos
	} else {
		fmt.Println("无缓存")
		rpcClient, err := Client.NewClientWithTarget(settings.Target.Addr)
		if err != nil {
			fmt.Println("rpc getTicketNumber failed, err:", err)
			return nil
		}

		response, err = rpcClient.GetTicketNumber(request)
		if response == nil {
			return nil
		}

		//缓存
		seatInfos := response.TrainsTicketInfo[0].SeatInfo
		for _, seatInfo := range seatInfos {
			RedisDB.HSet(cacheKey, strconv.Itoa(int(seatInfo.SeatTypeId)), seatInfo.SeatNumber)
		}
		RedisDB.Expire(cacheKey, time.Minute*1)
	}

	ticketInfos := response.TrainsTicketInfo

	if len(ticketInfos) == 0 {
		return nil
	}

	// 获取车次号

	ticketInfo := ticketInfos[0]
	train := &outer.Train{}
	trainNo := dao.GetTrainNumber(ticketInfo.TrainId)
	train.TrainNumber = trainNo
	train.TrainType = "G"
	train.TrainID = uint64(ticketInfo.TrainId)

	seatInfos := ticketInfo.SeatInfo
	for _, seatInfo := range seatInfos {
		switch seatInfo.SeatTypeId {
		case 0:
			train.BusinessSeat = int(seatInfo.SeatNumber)
			train.BusinessSeatPrice = 500
		case 1:
			train.FirstSeat = int(seatInfo.SeatNumber)
			train.FirstSeatPrice = 500
		case 2:
			train.SecondSeat = int(seatInfo.SeatNumber)
			train.SecondSeatPrice = 500
		case 3:
			train.SeniorSoftSleeper = int(seatInfo.SeatNumber)
			train.SeniorSoftBerthPrice = 500
		case 4:
			train.SoftSleeper = int(seatInfo.SeatNumber)
			train.SoftBerthPrice = 500
		case 5:
			train.HardSleeper = int(seatInfo.SeatNumber)
			train.HardBerthPrice = 500
		case 6:
			train.HardSeat = int(seatInfo.SeatNumber)
			train.HardSeatPrice = 500
		default:
		}
	}
	resMap, _ = RedisDB.HGetAll(trainNo).Result()
	leaveStation := dao.GetStationName(ssID)
	train.LeaveStation = leaveStation
	train.LeaveStationNo = uint64(ssID)
	train.LeaveTime, _ = RedisDB.HGet(trainNo+"-"+resMap[train.LeaveStation], "leaveTime").Result()
	arrivalStation := dao.GetStationName(esID)
	train.ArrivalStation = arrivalStation
	train.ArrivalStationNo = uint64(esID)
	train.ArrivalTime, _ = RedisDB.HGet(trainNo+"-"+resMap[train.ArrivalStation], "arriveTime").Result()

	if resMap[leaveStation] == "1" {
		train.LeaveStationType = "始"
	} else {
		train.LeaveStationType = "过"
	}
	if resMap[arrivalStation] == resMap["stationNum"] {
		train.ArrivalStationType = "终"
	} else {
		train.ArrivalStationType = "过"
	}

	// 获取始发站和终点站
	train.StartStation, _ = RedisDB.HGet(trainNo+"-"+"1", "stationName").Result()
	train.StartStationID, _ = RedisDB.HGet(trainNo+"-"+"1", "stationID").Result()
	train.EndStation, _ = RedisDB.HGet(trainNo+"-"+resMap["stationNum"], "stationName").Result()
	train.EndStationID, _ = RedisDB.HGet(trainNo+"-"+resMap["stationNum"], "stationID").Result()

	return train
}
