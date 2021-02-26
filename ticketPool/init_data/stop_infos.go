/*
* @Author: 余添能
* @Date:   2021/2/25 12:08 上午
 */
package init_data

import (
	"fmt"
	"time"
)

func InsertStopInfos()  {
	sqlStr:="insert into stop_infos(created_at,created_by,train_number,station_name,city,arrived_time,stay_duration,depart_time,stay_num,mileage)" +
		"values(?,?,?,?,?,?,?,?,?,?);"
	st,err:=Db.Prepare(sqlStr)
	if err!=nil{
		fmt.Println("Prepare table stop_infos failed, err:",err)
		return
	}
	trainNos:=ReadTotalTrainNo()
	for _,v:=range trainNos{
		for i:=0;i<int(v.StationNum);i++{
			station:=v.Stations[i]
			stayTime:=station.DepartTime.Sub(station.ArriveTime)
			stayDuration:=stayTime.Minutes()

			st.Exec(time.Now(),"系统",v.TrainNo,station.StationName,station.CityName,station.ArriveTime,int(stayDuration),
				station.DepartTime,i+1,station.Mileage)
		}

	}
}
