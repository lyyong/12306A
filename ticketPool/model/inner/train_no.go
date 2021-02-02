/*
* @Author: 余添能
* @Date:   2021/1/23 2:56 上午
 */
package inner

import "time"

//保存列车原始数据

type Train struct {
	ID           int64
	TrainNo      string     //车次
	StationNum   int64      //列车历经总站数
	Stations     []*Station //途经的所有站点
	InitialTime  time.Time  //起始时间
	TerminalTime time.Time  //停止时间
}
type Station struct {
	ID          int64
	StationNo   int64 //站序
	StationName string
	CityName    string
	ArriveTime  time.Time //到站时间
	DepartTime  time.Time //离站时间
	Duration    time.Time //从列车开始出发，距现在的时长
	Mileage     int64     //距离该站已行驶里程
	Price       float64   //从起始站开始的票价
}
