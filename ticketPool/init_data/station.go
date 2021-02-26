/*
* @Author: 余添能
* @Date:   2021/2/24 11:20 下午
 */
package init_data

import (
	"fmt"
	"time"
)

func WriteStations()  {
	sqlStr:="select city,station_name,station_spell from station_province_city;"
	rows,_:=Db.Query(sqlStr)
	sqlStr2:="insert into stations(created_at,created_by,name,city,spell,state) " +
		"values(?,?,?,?,?,?);"
	st,err:=Db.Prepare(sqlStr2)
	if err!=nil{
		fmt.Println("insert into stations failed,err:",err)
		return
	}
	for rows.Next(){
		var name,spell,city string
		rows.Scan(&city,&name,&spell)
		fmt.Println(city,name,spell)
		_,err:=st.Exec(time.Now(),"余添能",name,city,spell,1)
		if err!=nil{
			fmt.Println(err)
		}
	}
}


