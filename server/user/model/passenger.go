/**
 * @Author fzh
 * @Date 2021/2/22
 */
package model

import (
	"gorm.io/gorm"
)

type Passenger struct {
	gorm.Model
	Name              string
	CertificateType   int
	CertificateNumber string
	PassengerType     int
	State             int
	CreatedBy         string
	UpdatedBy         string
	DeletedBy         string
}

func InsertPassenger(db *gorm.DB, userId uint, passengers []*Passenger) error {
	user := new(User)
	user.ID = userId
	return db.Model(user).Association("Passengers").Append(passengers)
}

func UpdatePassenger(db *gorm.DB, userId uint, passengers []*Passenger) error {
	user := new(User)
	user.ID = userId
	return db.Model(user).Association("Passengers").Replace(passengers)
}

func ListPassenger(db *gorm.DB, userId uint) ([]Passenger, error) {
	user := new(User)
	user.ID = userId
	var passengers []Passenger
	err := db.Model(user).Association("Passengers").Find(&passengers)
	if err != nil {
		return nil, err
	}
	return passengers, nil
}
