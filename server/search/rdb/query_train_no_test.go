/*
* @Author: 余添能
* @Date:   2021/2/2 2:10 下午
 */
package rdb

import "testing"

func TestQueryTrainByCity(t *testing.T) {
	startCity := "北京"
	endCity := "上海"
	QueryTrainByCity(startCity, endCity)
}

func TestQueryTrainInfoByTrainNo(t *testing.T) {
	QueryTrainInfoByTrainNo("K5629", "合肥", "杭州")
}
