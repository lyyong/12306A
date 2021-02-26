/*
* @Author: 余添能
* @Date:   2021/2/26 12:07 上午
 */
package service

import (
	"fmt"
	"testing"
)

func TestQueryTrains(t *testing.T) {
	startCity:="北京市"
	endCity:="上海市"
	date:="2021-02-27"
	trains:=QueryTrains(startCity,endCity,date)
	fmt.Println(len(trains))
	for _,t:=range trains{
		fmt.Println(t)
	}
}
