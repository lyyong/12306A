/*
* @Author: 余添能
* @Date:   2021/2/23 10:08 下午
 */
package rdb

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
)

var (
	ConsulDb  *api.Client
	//Lockers map[string]*api.Lock

	 RedisDB *redis.Client
	 //ShaBuyTicket string
	 err error
)


func init()  {
	RedisDB =redis.NewClient(&redis.Options{
		Addr: "0.0.0.0:6379",
	})

	//加载脚本
	//buyTicketScript:= init_redis.CreateScriptBuyTicket()
	//ShaBuyTicket, err =buyTicketScript.Load(RedisDB).Result()
	//if err !=nil{
	//	fmt.Println("buyTicket lua script load failed ,err:", err)
	//	return
	//}
	//RedisDB = redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs:  []string{"192.168.10.11:7001","192.168.10.11:7002", "192.168.10.11:7003"},
	//})
	//连接redis集群


	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err == nil {
		ConsulDb = client
	}else{
		fmt.Println("connect consul failed, err:",err)
		return
	}

}
