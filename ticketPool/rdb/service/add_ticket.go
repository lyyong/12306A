/*
* @Author: 余添能
* @Date:   2021/2/24 12:22 上午
 */
package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"ticketPool/model/outer"
	"ticketPool/rdb"
)

func AddTicket(ticket *outer.Ticket) bool {

	////获取车次元数据
	//resMap,_:= rdb.RedisDB.HGetAll(ticket.TrainNumber).Result()
	startStationNo:=ticket.StartStationNum
	//fmt.Println(startStationNo)
	endStationNo,_:=strconv.ParseInt(ticket.EndStationNum,10,64)
	//fmt.Println(float64(endStationNo))
	dateAndTime:=strings.Split(ticket.StartTime," ")
	startTime:=dateAndTime[0]
	key:=startTime+":"+ticket.TrainNumber+":"+startStationNo+":"+ticket.SeatClass
	t,err:= rdb.RedisDB.ZAdd(key,redis.Z{Score: float64(endStationNo),Member: ticket.CarriageNum+":"+ticket.SeatNum}).Result()
	if err!=nil{
		fmt.Println("redis ZAdd failed, err:",err)
		return false
	}
	fmt.Println(t)
	return true
}
