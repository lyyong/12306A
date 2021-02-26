/*
* @Author: 余添能
* @Date:   2021/2/2 3:40 下午
 */
package rdb

import (
	"fmt"
	"testing"
)

func TestQueryTicketNumByTrainNoAndDate(t *testing.T) {
	QueryTicketNumByTrainNoAndDate("2021-1-23","K4729","secondSeat",1,10)
}

func TestQueryTrainByCity(t *testing.T) {
	date:="2021-02-26"
	startCity := "北京"
	endCity := "上海"
	trainNos := QueryTrainByDateAndCity(date, startCity, endCity)
	for _,t:=range trainNos{
		fmt.Println(t)
	}
	fmt.Println(len(trainNos))
}
