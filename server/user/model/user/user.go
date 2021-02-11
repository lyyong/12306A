/**
 * @Author fzh
 * @Date 2020/2/1
 */
package user

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
}

func InsertUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	user := new(User)
	result := db.Where("username = ?", username).First(&user)
	return user, result.Error
}
