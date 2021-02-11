/*
* @Author: 余添能
* @Date:   2021/2/4 10:43 下午
 */
package rdb

import (
	"fmt"
	"testing"
)

func TestQueryStation(t *testing.T) {
	QueryStation()
}

func TestQueryStationByTrainNo(t *testing.T) {
	stations := QueryStationByTrainNo("G21")
	for _,v:=range stations{
		fmt.Println(v)
	}
}