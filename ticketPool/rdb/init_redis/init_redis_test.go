/*
* @Author: 余添能
* @Date:   2021/1/31 7:16 下午
 */
package init_redis

import "testing"


//redis总初始化语句
func TestInitDataRedis(t *testing.T) {
	
	//静态数据
	WriteTrainPoolToRedis()
	WriteTrainInfoToRedis()
	WriteStationToRedis()
	WriteStationAndCityToRedis()

	//票池
	//WriteTicketPoolToRedis()
	//SplitTicket()
}



//=======================================
func TestWriteTrainPoolToRedis(t *testing.T) {
	WriteTrainPoolToRedis()
}

func TestWriteStationToRedis(t *testing.T) {
	WriteStationToRedis()
}
func TestWriteTrainInfoToRedis(t *testing.T) {
	WriteTrainInfoToRedis()
}


func TestWriteStationAndCityToRedis(t *testing.T) {
	WriteStationAndCityToRedis()
}
//北京市-上海市 0 1000
