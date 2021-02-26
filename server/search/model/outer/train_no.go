/*
* @Author: 余添能
* @Date:   2021/2/2 2:35 下午
 */
package outer


type Train struct {
	TrainNumber          string `json:"train_number"`

	StartTime        string `json:"leave_time"`
	//StartStationNo   int64  `json:"leave_station_no"`
	StartStation     string `json:"leave_station"`
	StartStationType string `json:"leave_station_type"`

	EndTime          string `json:"arrival_time"`
	//EndStationNo     int64  `json:"end_station_no"`
	EndStation       string `json:"arrival_station"`
	EndStationType   string `json:"arrival_station_type"`
	Duration         string `json:"travel_time"`

	TrainType        string `json:"train_type"`
	//高铁
	SecondSeat   int `json:"second_seats_number"`
	FirstSeat    int `json:"first_seats_number"`
	BusinessSeat int `json:"business_seats_number"`

	//火车
	SoftSleeper       int `json:"soft_berth_number"`
	HardSleeper       int `json:"hard_berth_number"`
	HardSeat          int `json:"hard_seats_number"`
	NoSeat            int `json:"no_seats_number"`
	SeniorSoftSleeper int `json:"senior_soft_berth_number"`
}
