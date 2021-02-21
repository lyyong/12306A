/*
* @Author: 余添能
* @Date:   2021/2/20 12:27 下午
 */
package rdb

import (
	"12306A/ticketPool/model/outer"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
)

//时间和日期都需要，不能显示已经离站的车次，买一张票，不选座
func BuyTicketOneNoSelect(buyTicket *outer.BuyTicket) *outer.Ticket {
	slices:=strings.Split(buyTicket.Date," ")
	//获得日期，不含时间，如：2021-02-01
	date:=slices[0]
	resMap,_:= RedisDB.HGetAll(buyTicket.TrainNo).Result()
	startStationNo,_:=resMap[buyTicket.StartStation]
	key := date+":"+buyTicket.TrainNo+":"+startStationNo+":"+buyTicket.SeatClass

	endStationNo:=resMap[buyTicket.EndStation]
	//RedisDB.ZRangeByScoreWithScores(key,redis.ZRangeBy{Min: endStationNo,Max: "1000"})
	//lua+redis
	res,err:=RedisDB.EvalSha(shaBuyTicket,[]string{key},endStationNo).Result()

	if err!=nil{
		fmt.Println("buy ticket by lua script failed, err:",err)
		return nil
	}
	//fmt.Println(res)
	result:=res.([]interface{})
	if strings.Compare(endStationNo,result[0].(string))!=0{
		//余票写回
		remainderEndStationNo,_:=strconv.ParseFloat(result[0].(string),64)
		remainderKey:=date+":"+buyTicket.TrainNo+":"+endStationNo+":"+buyTicket.SeatClass
		RedisDB.ZAdd(remainderKey,redis.Z{Score: remainderEndStationNo,Member: result[1].(string)})
		//fmt.Println(remainderRes)
	}
	//result[1]=车厢号:座位号
	carriageAndSeatNo:=strings.Split(result[1].(string),":")
	carriageNo:=carriageAndSeatNo[0]
	seatNo:=carriageAndSeatNo[1]

	ticket:=&outer.Ticket{}
	ticket.Date=buyTicket.Date
	ticket.TrainNo=buyTicket.TrainNo
	ticket.StartStation=buyTicket.StartStation
	ticket.EndStation=buyTicket.EndStation
	ticket.SeatClass=buyTicket.SeatClass
	ticket.CarriageNo=carriageNo
	ticket.SeatNo=seatNo
	//fmt.Println(ticket)
	return ticket
}


//买多张票
func BuyTicketMoreNoSelect(buyTicket *outer.BuyTicket)  {
	slices:=strings.Split(buyTicket.Date," ")
	//获得日期，不含时间，如：2021-02-01
	date:=slices[0]
	resMap,_:= RedisDB.HGetAll(buyTicket.TrainNo).Result()
	startStationNo,_:=resMap[buyTicket.StartStation]
	key := date+":"+buyTicket.TrainNo+":"+startStationNo+":"+buyTicket.SeatClass

	endStationNo:=resMap[buyTicket.EndStation]
}


