package ticket

import (
	"gorm.io/gorm"
	"time"
)

type Ticket struct{
	gorm.Model

	TrainId int32
	StartStation string
	DestStation string
	StartTime time.Time
	SeatType string
	CarriageNumber string
	SeatNumber string
	Amount string
	IndentId int32
	PassengerId int32
	State int8
}

func AddMultipleTicket(db *gorm.DB, tickets *[]Ticket) error {
	res := db.Create(tickets)
	return res.Error
}

func IsConflict(db *gorm.DB, passengerId *[]int32, date string) (bool, error) {
	return false, nil
}