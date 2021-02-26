/*
* @Author: 余添能
* @Date:   2021/1/30 11:23 下午
 */
package extra

type TotalTrainNo struct {
	ID           int
	TrainNo      string //车次
	StationNum   int64  //列车历经总站数
	InitialTime  string //起始时间
	TerminalTime string //停止时间
	StationNo    int64  //站序
	StationName  string
	CityName     string
	ArriveTime   string  //到站时间
	DepartTime   string  //离站时间
	Duration     string  //从列车开始出发，距现在的时长
	Mileage      int64   //距离该站已行驶里程
	Price        float64 //从起始站开始的票价
}
