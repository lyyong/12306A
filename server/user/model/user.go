/**
 * @Author fzh
 * @Date 2020/2/1
 */
package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	CreatedBy         string
	UpdatedBy         string
	DeletedBy         string
	Username          string
	Password          string
	State             int
	Salt              string
	CertificateType   int
	Name              string
	CertificateNumber string
	PhoneNumber       string
	Email             string
	PassengerType     int
	Passengers        []Passenger `gorm:"many2many:user_passengers;"`
}

func InsertUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUserById(db *gorm.DB, id uint) (*User, error) {
	user := new(User)
	result := db.First(&user, id)
	return user, result.Error
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	user := new(User)
	result := db.Where("username = ?", username).First(&user)
	return user, result.Error
}
