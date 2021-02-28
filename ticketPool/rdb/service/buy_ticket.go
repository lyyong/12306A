/*
* @Author: 余添能
* @Date:   2021/2/20 12:27 下午
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

//时间和日期都需要，不能显示已经离站的车次，买一张票，不选座
func BuyTicket(buyTickets []*outer.BuyTicket) (tickets []*outer.Ticket) {
	buyTicket:=buyTickets[0]
	slices:=strings.Split(buyTicket.StartTime," ")
	//获得日期，不含时间，如：2021-02-01
	date:=slices[0]
	//获取车次元数据
	resMap,_:= rdb.RedisDB.HGetAll(buyTicket.TrainNumber).Result()
	//获取上车站的站序
	startStationNo,_:=resMap[buyTicket.StartStation]

	key := date+":"+buyTicket.TrainNumber+":"+startStationNo+":"+buyTicket.SeatClass

	//使用consul提供的分布式锁
	lock,err:= rdb.ConsulDb.LockKey(key)
	_,err=lock.Lock(nil)
	if err!=nil{
		fmt.Println("get lock failed, err:",err)
		return nil
	}
	//不选座
	if buyTicket.SeatPlace==""{
		tickets= BuyTicketNoSelect(key,resMap,buyTickets)
	}else{
		tickets= BuyTicketNoSelect(key,resMap,buyTickets)
	}
	lock.Unlock()

	return tickets
}

func BuyTicketNoSelect(key string,trainInfo map[string]string,buyTickets[]*outer.BuyTicket) (tickets []*outer.Ticket) {
	if buyTickets == nil {
		return nil
	}
	ticketNum := len(buyTickets)
	startStationNo:=trainInfo[buyTickets[0].StartStation]
	endStationNo := trainInfo[buyTickets[0].EndStation]
	res, err := rdb.RedisDB.ZRangeByScoreWithScores(key, redis.ZRangeBy{Min: endStationNo, Max: "1000"}).Result()
	if err != nil {
		fmt.Println("zrangebyscore failed, err:", err)
		return nil
	}
	if len(res) < ticketNum {
		fmt.Println("无票")
		return nil
	}

	for i := 0; i < ticketNum; i++ {
		ticket := &outer.Ticket{}
		ticket.TrainNumber = buyTickets[i].TrainNumber
		ticket.StartTime = buyTickets[i].StartTime
		ticket.StartStation = buyTickets[i].StartStation
		ticket.StartStationNum=startStationNo
		ticket.EndTime = buyTickets[i].EndTime
		ticket.EndStation = buyTickets[i].EndStation
		ticket.EndStationNum=endStationNo
		ticket.SeatClass = buyTickets[i].SeatClass
		//res[1]=车厢号:座位号
		carriageAndSeatNo := strings.Split(res[i].Member.(string), ":")
		carriageNo := carriageAndSeatNo[0]
		seatNo := carriageAndSeatNo[1]
		ticket.CarriageNum = carriageNo
		ticket.SeatNum = seatNo
		tickets=append(tickets,ticket)
		//删掉原票
		rdb.RedisDB.ZRem(key,res[i].Member)
		//写回余票
		remainderStartStationNo := endStationNo
		//fmt.Println(remainderStartStationNo)
		remainderEndStationNo := strconv.Itoa(int(res[i].Score))
		//fmt.Println(res[i])
		//fmt.Println(remainderStartStationNo,remainderEndStationNo)
		//没有余票
		if strings.Compare(remainderStartStationNo, remainderEndStationNo) == 0 {
			//fmt.Println("没有余票")
			continue
		}
		////获得日期，不含时间，如：2021-02-01
		//stationMap,_:=rdb.RedisDB.HGetAll(buyTickets[i].TrainNumber+"-"+remainderStartStationNo).Result()
		startTime := strings.Split(ticket.EndTime," ")[0]
		fmt.Println(startTime)
		remainderKey := startTime + ":" + buyTickets[i].TrainNumber + ":" + remainderStartStationNo + ":" + buyTickets[i].SeatClass
		fmt.Println(remainderKey)
		r,err:=rdb.RedisDB.ZAdd(remainderKey, res[i]).Result()
		fmt.Println(r,err)

	}
	return tickets
}

//还未完全实现
func BuyTicketSelect(key string,trainInfo map[string]string,buyTickets[]*outer.BuyTicket) (tickets []*outer.Ticket) {
	if buyTickets == nil {
		return nil
	}
	ticketNum := len(buyTickets)
	startStationNo:=trainInfo[buyTickets[0].StartStation]
	endStationNo := trainInfo[buyTickets[0].EndStation]
	res, err := rdb.RedisDB.ZRangeByScoreWithScores(key, redis.ZRangeBy{Min: endStationNo, Max: "1000"}).Result()
	if err != nil {
		fmt.Println("zrangebyscore failed, err:", err)
		return nil
	}
	if len(res) < ticketNum {
		fmt.Println("无票")
		return nil
	}


	for i := 0; i < ticketNum; i++ {
		ticket := &outer.Ticket{}
		ticket.TrainNumber = buyTickets[i].TrainNumber
		ticket.StartTime = buyTickets[i].StartTime
		ticket.StartStation = buyTickets[i].StartStation
		ticket.StartStationNum=startStationNo
		ticket.EndTime = buyTickets[i].EndTime
		ticket.EndStation = buyTickets[i].EndStation
		ticket.EndStationNum=endStationNo
		ticket.SeatClass = buyTickets[i].SeatClass
		//res[1]=车厢号:座位号
		carriageAndSeatNo := strings.Split(res[i].Member.(string), ":")
		carriageNo := carriageAndSeatNo[0]
		seatNo := carriageAndSeatNo[1]
		ticket.CarriageNum = carriageNo
		ticket.SeatNum = seatNo
		tickets=append(tickets,ticket)

		//余票
		remainderStartStationNo := endStationNo
		remainderEndStationNo := strconv.Itoa(int(res[i].Score))
		//没有余票
		if strings.Compare(remainderStartStationNo, remainderEndStationNo) == 0 {
			continue
		}
		//获得日期，不含时间，如：2021-02-01
		startTime := trainInfo["departTime"]
		remainderKey := startTime + ":" + buyTickets[i].TrainNumber + ":" + remainderStartStationNo + ":" + buyTickets[i].SeatClass
		rdb.RedisDB.ZAdd(remainderKey, res[i])

	}
	return tickets
}
