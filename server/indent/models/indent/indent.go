package indent

import (
	"gorm.io/gorm"
	"time"
)

type Indent struct{
	gorm.Model
	IndentOuterId string
	UserId int32
	TrainId int32
	StartStation string
	DestStation string
	StartTime time.Time
	TicketNumber int8
	Amount string
	State int8
}


func AddIndent (db *gorm.DB, indent *Indent) error {
	return db.Create(indent).Error

}

func GetAllIndent (db *gorm.DB, userId int) ([]*Indent, error) {
	var indents []*Indent
	err := db.Where("user_id = ?", userId).Find(&indents).Error
	return indents, err
}

func GetNotStartIndent (db *gorm.DB, userId int) ([]*Indent, error) {
	var indents []*Indent
	err := db.Where("user_id = ? AND start_time > ?", userId, time.Now()).Find(&indents).Error
	return indents, err
}