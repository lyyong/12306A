/*
* @Author: 余添能
* @Date:   2021/2/25 6:11 下午
 */
package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"ticketPool/rdb"
	"time"
)

//date="2021-02-25
func QueryTrains(startCity,endCity string,date string) []string {
	key:=startCity+"-"+endCity
	now:=time.Now()
	nn:=now.Format("2006-01-02")
	min:=0
	if strings.Compare(date,nn)==0{
		//当天
		h,m,_:=now.Clock()
		min=h*60+m
	}

	res,err := rdb.RedisDB.ZRangeByScore(key, redis.ZRangeBy{Min: strconv.Itoa(min), Max: "50000"}).Result()
	if err!=nil{
		fmt.Println("select trains failed, err:",err)
		return nil
	}
	return res
}
