/*
* @Author: 余添能
* @Date:   2021/2/25 1:34 上午
 */
package dao

import (
	"12306A-search/model/inner"
	"fmt"
)

func SelectStopInfoAll() (stopInfos []*inner.StopInfo) {
	//created_at,updated_at,deleted_at,created_by,updated_by,deleted_by,
	sqlStr := "select id,train_id,station_id,train_number,station_name,city,arrived_time,leave_time,stop_seq " +
		"from stop_infos;"

	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		stopInfo := &inner.StopInfo{}
		//&stopInfo.CreatedAt,&stopInfo.UpdatedAt,&stopInfo.DeletedAt,&stopInfo.CreatedBy,&stopInfo.UpdatedBy,&stopInfo.DeletedBy,
		rows.Scan(&stopInfo.ID, &stopInfo.TrainId, &stopInfo.StationId, &stopInfo.TrainNumber, &stopInfo.StationName, &stopInfo.City, &stopInfo.ArrivedTime, &stopInfo.LeaveTime, &stopInfo.StopSeq)
		//fmt.Println(stopInfo.StationName)
		stopInfos = append(stopInfos, stopInfo)
	}

	return stopInfos
}
