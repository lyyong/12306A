/*
* @Author: 余添能
* @Date:   2021/2/24 11:46 下午
 */
package init_data

import (
	"fmt"
	"time"
)

func InsertTrains()  {
	trainNos:=ReadTotalTrainNo()
	sqlStr:="insert into trains(created_at,created_by,number,start_station,end_station,train_type,state) values(?,?,?,?,?,?,?);"
	st,err:=Db.Prepare(sqlStr)
	if err!=nil{
		fmt.Println("prepare insert trains failed, err:",err)
		return
	}
	for _,v:=range trainNos{
		_,err:=st.Exec(time.Now(),"余添能",v.TrainNo,v.Stations[0].StationName,v.Stations[v.StationNum-1].StationName,1,1)
		if err!=nil{
			fmt.Println("exec insert trains failed, err:",err)
		}
	}
}
