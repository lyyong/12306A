/*
* @Author: 余添能
* @Date:   2021/2/25 1:35 上午
 */
package dao

import (
	"12306A-search/model/inner"
	"fmt"

)

func SelectStationAll() []*inner.Station {
	sqlStr:="select id,city,name,spell from stations"
	rows,err:=Db.Query(sqlStr)
	if err!=nil{
		fmt.Println("select stations failed, err:",err)
		return nil
	}
	var stations []*inner.Station
	for rows.Next(){
		station:=&inner.Station{}
		rows.Scan(&station.ID,&station.City,&station.Name,&station.Spell)
		stations=append(stations,station)
	}
	return stations
}
