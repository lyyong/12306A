package models

import (
	"gorm.io/gorm"
	"ticket/utils/database"
	"time"
)

type Ticket struct {
	gorm.Model

	UserId         uint32
	TrainId        uint32
	TrainNum       string
	StartStationId uint32
	StartStation   string
	StartTime      time.Time
	DestStationId  uint32
	DestStation    string
	DestTime       time.Time
	SeatType       string
	CarriageNumber string
	SeatNumber     string
	Price          int32
	OrderOutsideId string
	PassengerName  string
	PassengerId    uint32
	State          uint8
}

const (
	TicketBuySuccessful = iota
	TicketFinish
	TicketRefund
	TicketChange
	TicketWaitCash
	TicketRefundFinish
)

func AddMultipleTicket(tickets *[]Ticket) error {
	res := database.DB.Create(tickets)
	return res.Error
}

func GetTicketByOrderId(orderId string) ([]*Ticket, error) {
	var tickets []*Ticket
	res := database.DB.Where("order_outside_id = ?", orderId).Find(&tickets)
	return tickets, res.Error
}

func GetTicketsByPassengerId(passengerId uint32) ([]*Ticket, error) {
	var tickets []*Ticket
	res := database.DB.Where("passenger_id = ?", passengerId).Find(&tickets)
	return tickets, res.Error
}

func GetTicketByID(id uint) (*Ticket, error) {
	var ticket Ticket
	err := database.DB.Where("id = ?", id).Find(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func DeleteTicketByTicketId(db *gorm.DB, ticketsId []uint32) ([]*Ticket, error) {
	var tickets []*Ticket
	res := db.Where("id IN ?", ticketsId).Find(&tickets)
	if res.Error != nil {
		return nil, res.Error
	}
	res = db.Table("tickets").Where("id IN ?", ticketsId).Updates(Ticket{State: 2})
	return tickets, res.Error
}

func UpdateState(ticketsId []uint32, state uint8) error {
	res := database.DB.Table("tickets").Where("id IN ?", ticketsId).Updates(Ticket{State: state})
	return res.Error
}

func IsConflict(passengerId *[]uint32, date string) (bool, error) {
	return false, nil
}
