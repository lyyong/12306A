// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package model

import (
	"common/tools/logging"
	"gorm.io/gorm"
	"ticketPool/utils/database"
)

type StopInfo struct {
	Model
	TrainID     uint   // 火车车次id
	StationID   uint   // 站点id
	TrainNumber string // 车次编号
	StationName string // 车站名称
	City        string // 车站所在的城市
	ArrivedTime string // 到达时间
	LeaveTime   string // 离开时间
	StopSeq     int    // 停靠顺序
}

// GetStopInfoByTrainID 通过trainID获取停靠站信息
func GetStopInfoByTrainID(trainID uint) []*StopInfo {
	sis := make([]*StopInfo, 0)
	err := database.DB.Where("train_id = ?", trainID).Find(&sis).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(err)
		return nil
	}
	return sis
}
