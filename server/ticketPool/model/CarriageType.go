// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package model

import (
	"common/tools/logging"
	"gorm.io/gorm"
	"ticketPool/utils/database"
)

type CarriageType struct {
	Model
	// 座位数量
	SoftBerthNumber       int
	HardBerthNumber       int
	SeniorSoftBerthNumber int
	HardSeatNumber        int
	SecondSeatNumber      int
	FirstSeatNumber       int
	BusinessSeatNumber    int
	// 座位编号
	BusinessSeat    string
	FirstSeat       string
	SecondSeat      string
	HardSeat        string
	HardBerth       string
	SoftBerth       string
	SeniorSoftBerth string
}

// GetCarriageTypes 获取所有的车厢类型, 出错返回nil
func GetCarriageTypes() []*CarriageType {
	carriageTypes := make([]*CarriageType, 0)
	if err := database.DB.Find(&carriageTypes).Error; err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return nil
	}
	return carriageTypes
}

// GetCarriageTypesByID 通过id获取车厢类型, 出错返回nil
func GetCarriageTypesByID(id uint) *CarriageType {
	var ct *CarriageType
	if err := database.DB.Find(ct).Where("id = ?", id).Error; err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return nil
	}
	return ct
}
