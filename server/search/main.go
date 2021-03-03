/*
* @Author: 余添能
* @Date:   2021/2/4 10:57 下午
 */
package main

import (
	"12306A-search/dao"
	"12306A-search/rdb"
	"12306A-search/router"
	"12306A-search/tools/settings"
	"common/tools/logging"
	"fmt"
)

func main()  {
	//fmt.Println("aa"=="aa")
	logging.Info("加载配置文件")
	settings.Setup()
	logging.Info("初始化DB")
	dao.InitDB()
	logging.Info("初始化redis")
	rdb.InitRedis()
	logging.Info("启动服务器...")
	r:=router.InitRouter()
	fmt.Println(r)

}
//北京市-上海市 800 10000000