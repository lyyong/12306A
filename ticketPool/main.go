/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"common/tools/logging"
	"ticketPool/rdb/init_redis"
	"ticketPool/rpc"
	"ticketPool/tools/setting"
)


func main()  {
	//dao.InitId()
	init_redis.InitDataRedis()
	//fmt.Println("初始化完成")
	logging.Info("初始化完成")
	setting.Setup()
	rpc.Setup()

}


