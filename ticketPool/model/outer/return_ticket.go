/*
* @Author: 余添能
* @Date:   2021/2/4 10:01 下午
 */
package outer

//已购票信息
type Ticket struct {
	ID 				int
	Date           string
	TrainNumber        string
	StartTime      string
	StartStationNum string
	StartStation   string
	EndTime        string
	EndStation     string
	EndStationNum   string
	SeatClass      string
	CarriageNum     string
	SeatNum         string
	Price          float64
}
