/*
* @Author: 余添能
* @Date:   2021/2/4 10:33 下午
 */
package rdb

import "fmt"

//返回所有站点
//用一个list保存，key=stations  站点
func QueryStation() []string {
	key:="stations"
	stations,err:=RedisDB.LRange(key,0,10000).Result()
	if err!=nil{
		fmt.Println("query stations failed, err:",err)
		return nil
	}
	for _,v:=range stations{
		fmt.Println(v)
	}
	fmt.Println(len(stations))
	return stations
}
