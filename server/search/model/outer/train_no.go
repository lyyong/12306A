/*
* @Author: 余添能
* @Date:   2021/2/2 2:35 下午
 */
package outer

type Train struct {
	TrainNo      string
	StartTime    string
	StartStationNo int64
	StartStation string
	EndTime      string
	EndStationNo int64
	EndStation   string
	Duration     string
	//高铁
	SecondSeat   int
	FirstSeat    int
	BusinessSeat int

	//火车
	SoftSleeper int
	HardSleeper int
	HardSeat    int
}
