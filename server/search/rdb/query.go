/*
* @Author: 余添能
* @Date:   2021/1/26 3:37 下午
 */
package rdb

import (
	"fmt"
	"github.com/go-redis/redis"
)

var RedisDB *redis.Client

func init() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func QueryStation() {
	key := "stationCity"
	resMap, _ := RedisDB.HGetAll(key).Result()
	for k, v := range resMap {
		fmt.Println(k, v)
	}
}
