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
	PhoneNumber       string
	PassengerType     int
	State             int
	CreatedBy         string
	UpdatedBy         string
	DeletedBy         string
}

func InsertPassenger(db *gorm.DB, userId uint, passenger *Passenger) error {
	user := new(User)
	user.ID = userId
	return db.Model(user).Association("Passengers").Append(passenger)
}

func UpdatePassenger(db *gorm.DB, userId uint, passenger *Passenger) error {
	return db.Model(passenger).Updates(passenger).Error
}

func DeletePassenger(db *gorm.DB, userId uint, passenger *Passenger) error {
	user := new(User)
	user.ID = userId
	return db.Model(user).Association("Passengers").Delete(passenger)
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
