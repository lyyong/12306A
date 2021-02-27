/*
* @Author: 余添能
* @Date:   2021/2/25 1:34 上午
 */
package dao

import (
	"fmt"
	"ticketPool/model/inner"
)

func SelectTrainAll() (trains []*inner.Train) {
	sqlStr:="select id,number,start_station,end_station,train_type,state from trains;"
	rows,err:=Db.Query(sqlStr)
	if err!=nil{
		fmt.Println("select table trains failed, err:",err)
		return nil
	}
	for rows.Next(){
		train:=&inner.Train{}
		rows.Scan(&train.ID,&train.Number,&train.StartStation,&train.EndStation,&train.TrainType,&train.State)
		//fmt.Println(train)
		trains=append(trains,train)
	}
	return trains
}
