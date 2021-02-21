// @Author LiuYong
// @Created at 2021-02-21
// @Modified at 2021-02-21
package model

import (
	"gorm.io/gorm"
	"ticketPool/utils/database"
)

type Train struct {
	Model
	Number       string
	StartStation string
	EndStation   string
	TrainType    uint
	State        int
}

// GetTrainsByCondition 通过条件获取车次
func GetTrainsByCondition(condition map[string]interface{}) []*Train {
	ts := make([]*Train, 0)
	err := database.DB.Where(condition).Find(&ts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return ts
}

// GetTrainsByNumberLike 通过车次编号模糊查找 %为通配符
func GetTrainsByNumberLike(like string) []*Train {
	ts := make([]*Train, 0)
	err := database.DB.Where("number LIKE ?", like).Find(&ts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return ts
}
