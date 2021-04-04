// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package model

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	CreatedBy string
	DeletedBy string
	UpdatedBy string
}
