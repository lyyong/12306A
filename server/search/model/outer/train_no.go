/*
* @Author: 余添能
* @Date:   2021/2/2 2:35 下午
 */
package outer

type Train struct {
	TrainNumber string `json:"train_number"`
	TrainID     uint64 `json:"train_id"`

	StartTime string `json:"leave_time"`

	StartStation     string `json:"leave_station"`
	StartStationNo   uint64 `json:"leave_station_id"`
	StartStationType string `json:"leave_station_type"`

	EndTime        string `json:"arrival_time"`
	EndStation     string `json:"arrival_station"`
	EndStationNo   uint64 `json:"arrival_station_id"`
	EndStationType string `json:"arrival_station_type"`
	Duration       string `json:"travel_time"`

	TrainType string `json:"train_type"`
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
