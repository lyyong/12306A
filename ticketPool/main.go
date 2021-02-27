/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"fmt"
	"ticketPool/rdb/init_redis"
)


func main()  {
	//dao.InitId()
	//rpc.Setup()
	init_redis.InitDataRedis()
	fmt.Println("初始化完成")
}


