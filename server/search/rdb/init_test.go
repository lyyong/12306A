/*
* @Author: 余添能
* @Date:   2021/5/29 4:00 下午
 */
package rdb

import "testing"

//zrangebyscore 北京市-上海市 0 1000 withscores
func TestWriteTrainPoolToRedis(t *testing.T) {
	WriteTrainPoolToRedis()

}

func TestWriteTrainInfoToRedis(t *testing.T) {
	WriteTrainInfoToRedis()
}

func TestWriteStationAndCityToRedis(t *testing.T) {
	WriteStationAndCityToRedis()
}
