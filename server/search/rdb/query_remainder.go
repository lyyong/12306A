/*
* @Author: 余添能
* @Date:   2021/1/26 3:37 下午
 */
package rdb

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

//查询某一天某个车次的某种座位等级的车票数，合适的
func QueryTicketNumByTrainNoAndDate(date,trainNo,seatClass string,start,end int) int64{
	//2021-1-23:K4729:1:secondSeat
	key := date+":"+trainNo+":"+strconv.Itoa(start)+":"+seatClass
	//fmt.Println(key)
	//寻找下车站大于等于end的票
	res,err:=RedisDB.ZRangeByScoreWithScores(key,redis.ZRangeBy{Min: strconv.Itoa(end),Max: "1000"}).Result()
	if err!=nil{
		fmt.Println("redis query zset failed, err:",err)
		return 0
	}
	return int64(len(res))
}


