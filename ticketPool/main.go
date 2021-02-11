/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"12306A/ticketPool/init_data"
	"12306A/ticketPool/rdb"
	"fmt"
)

func main()  {
	init_data.InitDataMysql()
	rdb.InitDataRedis()
	fmt.Println("初始化结束")

}
