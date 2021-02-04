/*
* @Author: 余添能
* @Date:   2021/2/3 7:12 下午
 */
package rdb

import (
	"12306A/server/search/model/outer"
	"strings"
)
//总查询
func Query(search *outer.Search) []*outer.Train {
	res:=strings.Split(search.Date," ")
	date:=res[0]

	//日期：2021-1-23

	var trains []*outer.Train
	//根据城市查车次
	trainNos := QueryTrainByCity(search.StartCity, search.EndCity)
	for _,t:=range trainNos{
		//根据车次查两站信息
		trainNo := QueryTrainInfoByTrainNo(t, search.StartCity, search.EndCity)

		//根据车次查询票数
		firstSeat:=QueryTicketNumByTrainNoAndDate(date,trainNo.TrainNo,"firstSeat",int(trainNo.StartStationNo),int(trainNo.EndStationNo))
		trainNo.FirstSeat=int(firstSeat)
		secondSeat:=QueryTicketNumByTrainNoAndDate(date,trainNo.TrainNo,"secondSeat",int(trainNo.StartStationNo),int(trainNo.EndStationNo))
		trainNo.SecondSeat=int(secondSeat)
		businessSeat:=QueryTicketNumByTrainNoAndDate(date,trainNo.TrainNo,"businessSeat",int(trainNo.StartStationNo),int(trainNo.EndStationNo))
		trainNo.BusinessSeat=int(businessSeat)
		//fmt.Println(trainNo)
		trains=append(trains,trainNo)
	}
	return trains
}
