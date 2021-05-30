/*
* @Author: 余添能
* @Date:   2021/2/2 2:35 下午
 */
package outer

type Train struct {
	TrainNumber string `json:"train_number"`
	TrainID     uint64 `json:"train_id"`

	LeaveTime        string `json:"leave_time"`
	LeaveStation     string `json:"leave_station"`
	LeaveStationNo   uint64 `json:"leave_station_id"`
	LeaveStationType string `json:"leave_station_type"`

	ArrivalTime        string `json:"arrival_time"`
	ArrivalStation     string `json:"arrival_station"`
	ArrivalStationNo   uint64 `json:"arrival_station_id"`
	ArrivalStationType string `json:"arrival_station_type"`

	Duration string `json:"travel_time"`

	// 列车的始发站和终点站
	StartStation   string `json:"start_station"`
	StartStationID string `json:"start_station_id"`
	EndStation     string `json:"end_station"`
	EndStationID   string `json:"end_station_id"`

	TrainType string `json:"train_type"`
	//高铁
	SecondSeat        int `json:"second_seats_number"`
	SecondSeatPrice   int `json:"second_seats_price"`
	FirstSeat         int `json:"first_seats_number"`
	FirstSeatPrice    int `json:"first_seats_price"`
	BusinessSeat      int `json:"business_seats_number"`
	BusinessSeatPrice int `json:"business_seats_price"`

	//火车
	SoftSleeper          int `json:"soft_berth_number"`
	SoftBerthPrice       int `json:"soft_berth_price"`
	HardSleeper          int `json:"hard_berth_number"`
	HardBerthPrice       int `json:"hard_berth_price"`
	HardSeat             int `json:"hard_seats_number"`
	HardSeatPrice        int `json:"hard_seats_price"`
	NoSeat               int `json:"no_seats_number"`
	NoSeatPrice          int `json:"no_seats_price"`
	SeniorSoftSleeper    int `json:"senior_soft_berth_number"`
	SeniorSoftBerthPrice int `json:"senior_soft_berth_price"`
}
