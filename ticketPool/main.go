/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"common/tools/logging"
	"ticketPool/dao"
	"ticketPool/rdb"
	"ticketPool/rdb/init_redis"
	"ticketPool/tools/setting"
)

func main()  {

	logging.Info("初始化配置文件")
	setting.Setup()
	logging.Info("初始化DB")
	dao.InitDB()
	logging.Info("初始化consul")
	rdb.InitConsul()
	logging.Info("初始化redis")
	rdb.InitRedis()
	logging.Info("初始化静态数据")
	init_redis.InitDataRedis()
	//logging.Info("初始化票池")
	//init_redis.InitTicketPool()
	//logging.Info("初始化rpc")
	//rpc.Setup()
}


