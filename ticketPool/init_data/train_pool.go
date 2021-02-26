/*
* @Author: 余添能
* @Date:   2021/1/31 12:03 上午
 */
package init_data

import (
	"fmt"
	"ticketPool/model/extra"
)


func ReadTrainPoolAll() []*extra.TrainPool {
	strSql := "select train_no,start_city,start_time,end_city,end_time from train_pools;"
	rows, err := Db.Query(strSql)
	if err != nil {
		fmt.Println("select train_pools failed, err:", err)
		return nil
	}
	var trainPools []*extra.TrainPool
	for rows.Next() {
		trainPool := &extra.TrainPool{}
		rows.Scan(&trainPool.TrainNo, &trainPool.StartCity, &trainPool.StartTime, &trainPool.EndCity, &trainPool.EndTime)
		trainPools = append(trainPools, trainPool)
	}
	fmt.Println(len(trainPools))
	return trainPools
}
