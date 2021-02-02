/*
* @Author: 余添能
* @Date:   2021/2/2 2:35 下午
 */
package outer

type Train struct {
	TrainNo      string
	StartTime    string
	StartStation string
	EndTime      string
	EndStation   string
	Duration     string
	//高铁
	Second   int
	First    int
	Business int

	//火车
	SoftSleeper int
	HardSleeper int
	HardSeat    int
}
