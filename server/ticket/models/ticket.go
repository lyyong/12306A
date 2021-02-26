package models

import (
	"gorm.io/gorm"
	"ticket/utils/database"
	"time"
)

type Ticket struct{
	gorm.Model

	UserId			uint32
	TrainId 		uint32
	TrainNum 		string
	StartStationId 	uint32
	StartStation 	string
	StartTime 		time.Time
	DestStationId 	uint32
	DestStation 	string
	DestTime 		time.Time
	SeatType 		string
	CarriageNumber 	string
	SeatNumber 		string
	Price 			int32
	OrderOutsideId 	string
	PassengerName 	string
	PassengerId 	uint32
	State 			uint8
}

func AddMultipleTicket(tickets *[]Ticket) error {
	res := database.DB.Create(tickets)
	return res.Error
}

func UpdateState(ticketId uint32, state string) (bool, error) {
	return false, nil
}

func IsConflict(passengerId *[]uint32, date string) (bool, error) {
	return false, nil
}