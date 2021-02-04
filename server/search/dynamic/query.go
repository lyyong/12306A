/*
* @Author: 余添能
* @Date:   2021/2/4 10:25 下午
 */
package dynamic

import (
	"12306A/server/search/model/outer"
	rdb "12306A/server/search/rdb"
	"fmt"
	"strings"
)

func Query(search *outer.Search) []*outer.Train {
	res:=strings.Split(search.Date," ")
	date:=res[0]

	//日期：2021-1-23

	var trains []*outer.Train
	//根据城市查车次
	trainNos := rdb.QueryTrainByCity(search.StartCity, search.EndCity)
	for _,t:=range trainNos{
		//根据车次查两站信息
		trainNo := rdb.QueryTrainInfoByTrainNo(t, search.StartCity, search.EndCity)

		//根据车次查询票数
		firstSeat:=rdb.QueryTicketNumByTrainNoAndDate(date,trainNo.TrainNo,"firstSeat",int(trainNo.StartStationNo),int(trainNo.EndStationNo))
		trainNo.FirstSeat=int(firstSeat)
		secondSeat:=rdb.QueryTicketNumByTrainNoAndDate(date,trainNo.TrainNo,"secondSeat",int(trainNo.StartStationNo),int(trainNo.EndStationNo))
		trainNo.SecondSeat=int(secondSeat)
		businessSeat:=rdb.QueryTicketNumByTrainNoAndDate(date,trainNo.TrainNo,"businessSeat",int(trainNo.StartStationNo),int(trainNo.EndStationNo))
		trainNo.BusinessSeat=int(businessSeat)
		fmt.Println(trainNo)
		trains=append(trains,trainNo)
	}
	return trains
}
