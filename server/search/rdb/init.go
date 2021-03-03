/*
* @Author: 余添能
* @Date:   2021/3/3 8:53 下午
 */
package rdb

import (
	"12306A-search/tools/settings"
	"fmt"
	"github.com/go-redis/redis"
)


var RedisDB *redis.Client
var shaBuyTicket string
var err error


func InitRedis()  {

	RedisDB = redis.NewClient(&redis.Options{
		Addr:     settings.RedisDB.Host,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	buyTicketScript:=CreateScriptBuyTicket()
	shaBuyTicket,err=buyTicketScript.Load(RedisDB).Result()
	if err!=nil{
		fmt.Println("buyTicket lua script load failed ,err:",err)
		return
	}
}
