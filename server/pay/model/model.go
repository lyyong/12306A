// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package model

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	CreatedBy string `json:"create_by"`
	UpdatedBy string `json:"updated_by"`
	DeletedBy string `json:"delete_by"`
}
