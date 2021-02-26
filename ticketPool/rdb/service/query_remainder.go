/*
* @Author: 余添能
* @Date:   2021/2/20 12:24 下午
 */
package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"ticketPool/rdb"
)

//查询某一天某个车次的某种座位等级的车票数，合适的
func QueryTicketNumByTrainNoAndDate(date,trainNo,seatClass string,startStation,endStation string) int64{
	//2021-02-25:G21:1:secondSeat
	resMap,_:=rdb.RedisDB.HGetAll(trainNo).Result()
	startStationNo:=resMap[startStation]
	endStationNo:=resMap[endStation]
	key := date+":"+trainNo+":"+startStationNo+":"+seatClass
	//fmt.Println(key)
	//寻找下车站大于等于end的票
	res,err:= rdb.RedisDB.ZRangeByScoreWithScores(key,redis.ZRangeBy{Min: endStationNo,Max: "10000"}).Result()
	if err!=nil{
		fmt.Println("redis query zset failed, err:",err)
		return 0
	}
	return int64(len(res))
}
