/*
* @Author: 余添能
* @Date:   2021/2/2 3:40 下午
 */
package rdb

import (
	"fmt"
	"testing"
)

func TestQueryTrainByCity(t *testing.T) {
	date:="2021-02-28"
	startCity := "北京市"
	endCity := "上海市"
	trainNos := QueryTrainByDateAndCity(date, startCity, endCity)
	for _,t:=range trainNos{
		fmt.Println(t)
	}
	fmt.Println(len(trainNos))
}
