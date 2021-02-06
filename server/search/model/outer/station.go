/*
* @Author: 余添能
* @Date:   2021/2/6 2:00 下午
 */
package outer

type Station struct {
	StationNo int `json:"station_no"`
	StationName string `json:"station_name"`
	ArriveTime string `json:"arrival_time"`
	DepartTime string `json:"leave_time"`
	WaitTime string `json:"wait_time"`
}
