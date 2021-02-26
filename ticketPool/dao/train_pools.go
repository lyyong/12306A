/*
* @Author: 余添能
* @Date:   2021/1/31 6:59 下午
 */
package dao

import (
	"fmt"
	"ticketPool/init_data"
	"ticketPool/model/inner"
)

func SelectTrainPoolAll() []*inner.TrainPool {
	strSql := "select id,train_no,start_city,start_time,end_city,end_time from train_pools;"

	rows, err := init_data.Db.Query(strSql)
	if err != nil {
		fmt.Println("select train_pools failed, err:", err)
		return nil
	}
	var trainPools []*inner.TrainPool
	for rows.Next() {
		trainPool := &inner.TrainPool{}
		//fmt.Println(trainPool)
		err:=rows.Scan(&trainPool.ID,&trainPool.TrainNo, &trainPool.StartCity,
			&trainPool.StartTime, &trainPool.EndCity, &trainPool.EndTime)
		if err!=nil{
			fmt.Println("scan failed")
		}
		trainPools = append(trainPools, trainPool)
		//fmt.Println(trainPool)
	}
	return trainPools
}

