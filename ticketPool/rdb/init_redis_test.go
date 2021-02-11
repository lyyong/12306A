/*
* @Author: 余添能
* @Date:   2021/1/31 7:16 下午
 */
package rdb

import "testing"

func TestWriteTrainPoolToRedis(t *testing.T) {
	WriteTrainPoolToRedis()
}

func TestWriteTrainInfoToRedis(t *testing.T) {
	WriteTrainInfoToRedis()
}
func TestWriteStationAndCityToRedis(t *testing.T) {
	WriteStationAndCityToRedis()
}

func TestWriteStationToRedis(t *testing.T) {
	WriteStationToRedis()
}
func TestWriteTicketPoolToRedis(t *testing.T) {
	WriteTicketPoolToRedis()
}