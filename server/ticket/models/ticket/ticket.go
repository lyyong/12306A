package ticket

import (
	"gorm.io/gorm"
	"time"
)

type Ticket struct{
	gorm.Model

	TrainId 		uint32
	TrainNum 		string
	StartStationId 	uint32
	StartStation 	string
	StartTime 		time.Time
	DestStationId 	uint32
	DestStation 	string
	DestTime 		time.Time
	SeatTypeId 		uint32
	SeatType 		string
	CarriageNumber 	string
	SeatNumber 		string
	Price 			uint32
	OrderId 		uint32
	PassengerName 	string
	PassengerId 	uint32
	State 			uint8
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