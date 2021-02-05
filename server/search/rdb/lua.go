/*
* @Author: 余添能
* @Date:   2021/2/4 6:21 下午
 */
package rdb

import (
	"fmt"
	"github.com/go-redis/redis"
)


var RedisDB *redis.Client
var shaBuyTicket string
var err error

func init() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
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

func CreateScriptBuyTicket() *redis.Script {
	script := redis.NewScript(`
		local key=tostring(KEYS[1])
		local min=tonumber(ARGV[1])
		local tickets = redis.call("ZRangeByScore", key, min,1000,"withscores")

		-- 表的大小为0，表示没有元素
		if #tickets==0 then 
		   return {0,0,0}
		end
		
		-- 奇数key 对应val是　value
    	-- 偶数key 对应val是　score

		local endStation
		local carriageAndSeatNo
		for i,v in pairs(tickets) do
			if i==1 then
				carriageAndSeatNo=v
			end
			if i==2 then
				endStation=v
				break;
			end
		end
		print(carriageAndSeatNo)
		local res=redis.call("ZREM",key,carriageAndSeatNo)
		return {endStation,carriageAndSeatNo,res}
	`)
	return script
}

