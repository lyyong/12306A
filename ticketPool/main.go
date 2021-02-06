/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func main()  {
	//init_data.InitDataMysql()
	//rdb.InitDataRedis()

	RedisDB := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{ // 填写master主机
			"192.168.10.11:7001",
			"192.168.10.11:7002",
			"192.168.10.11:7003",
		},
		RouteRandomly: true,
		DialTimeout:  5000 * time.Microsecond, // 设置连接超时
		ReadTimeout:  5000 * time.Microsecond, // 设置读取超时
		WriteTimeout: 5000 * time.Microsecond, // 设置写入超时
	})
	//RedisDB:=redis.NewClient(&redis.Options{
	//	Addr: "192.168.10.11:7002",
	//})
	// 发送一个ping命令,测试是否通
	s,err:= RedisDB.Do("ping").String()
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(s)
	RedisDB.Set("eee",111,time.Second*100)
	RedisDB.Set("ppp",222,time.Second)
	fmt.Println(RedisDB.Get("ppp").Result())
	fmt.Println(RedisDB.Get("b").Result())
	fmt.Println(RedisDB.Get("eee"))
}
