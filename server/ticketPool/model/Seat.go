// @Author: KongLingWen
// @Created at 2021/2/22
// @Modified at 2021/2/22

package model

import (
	"ticketPool/utils/database"
)

type Seat struct {
	Model
	ID         uint
	TrainId    uint32
	Date       string
	SeatTypeId uint32
	Key        uint64
	SeatInfo   string
}

func DeleteSeat(seat *Seat, value []string) error {
	return database.DB.Where("train_id = ? and date = ? and seat_type_id = ? and seat_info in ?", seat.TrainId, seat.Date, seat.SeatTypeId, value).Delete(&Seat{}).Error
}

func InsertSeat(seats *[]Seat) error {
	return database.DB.Create(seats).Error
}
