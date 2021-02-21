// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package model

import (
	"common/tools/logging"
	"gorm.io/gorm"
	"ticketPool/utils/database"
)

type TrainType struct {
	Model
	TypeNumber    string // 列车型号
	CarriageList  string // 车厢列表
	CarriageNum   int    // 车厢数
	MaxSpeed      int    // 最大速度
	WifiState     int    // wifi状态
	FoolCarriages string // 餐车列表
	MaxPassenger  int    // 最大再员
	Length        int    // 列车长度
}

// GetTrainTypeByID 通过id获取火车类型
func GetTrainTypeByID(id uint) *TrainType {
	var tt *TrainType
	if err := database.DB.Find(tt).Where("id = ?", id).Error; err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return nil
	}
	return tt
}

// GetTrainTypes 获取有所列车类型
func GetTrainTypes() []*TrainType {
	tts := make([]*TrainType, 0)
	err := database.DB.Find(&tts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return nil
	}
	return tts
}
