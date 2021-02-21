package ticket

import (
	"gorm.io/gorm"
	"time"
)

type Ticket struct{
	gorm.Model

	TrainId uint32
	StartStation string
	DestStation string
	StartTime time.Time
	SeatType string
	CarriageNumber string
	SeatNumber string
	Amount string
	IndentId uint32
	PassengerId uint32
	State int8
}

func AddMultipleTicket(db *gorm.DB, tickets *[]Ticket) error {
	res := db.Create(tickets)
	return res.Error
}

func UpdateState(db *gorm.DB, ticketId uint32, state string) (bool, error) {
	return false, nil
}

func IsConflict(db *gorm.DB, passengerId *[]uint32, date string) (bool, error) {
	return false, nil
}