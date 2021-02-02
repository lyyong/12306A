/*
* @Author: 余添能
* @Date:   2021/2/2 3:10 下午
 */
package service

import (
	outer2 "12306A/server/search/model/outer"
	rdb2 "12306A/server/search/rdb"
	"fmt"
)

func QueryTrainByCity(startCity, endCity string) []*outer2.Train {

	trainNos := rdb2.QueryTrainByCity(startCity, endCity)

	var trains []*outer2.Train
	for _, t := range trainNos {
		res := rdb2.QueryTrainInfoByTrainNo(t, startCity, endCity)
		trains = append(trains, res)
	}

	for _, v := range trains {
		fmt.Println(v)
	}
	return trains
}
