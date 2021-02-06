/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"12306A/ticketPool/init_data"
	"12306A/ticketPool/rdb"
)

func main()  {
	init_data.InitDataMysql()
	rdb.InitDataRedis()

}
