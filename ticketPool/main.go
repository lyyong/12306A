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

	//RedisDB := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: []string{ // 填写master主机
	//		"192.168.10.11:7001",
	//		"192.168.10.12:7002",
	//		"192.168.10.13:7003",
	//	},
	//	DialTimeout:  50 * time.Microsecond, // 设置连接超时
	//	ReadTimeout:  50 * time.Microsecond, // 设置读取超时
	//	WriteTimeout: 50 * time.Microsecond, // 设置写入超时
	//})
	//// 发送一个ping命令,测试是否通
	//s,err:= RedisDB.Do("ping").String()
	//if err!=nil{
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(s)
}
