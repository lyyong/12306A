// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package model

import (
	"pay/model"
	"pay/tools/database"
	"time"

	"gorm.io/gorm"
)

type Candidate struct {
	model.Model
	Date           time.Time `json:"date"`             // 候补时间
	TrainID        uint      `json:"train_id"`         // 车次id
	OrderID        string    `json:"order_id"`         // 订单id
	UserID         uint      `json:"user_id"`          // 用户id
	PassengerID    uint      `json:"passenger_id"`     // 乘客id
	PassengerName  string    `json:"passenger_name"`   // 乘客名称
	StartStationID uint      `json:"start_station_id"` // 上车车站id
	DestStationID  uint      `json:"dest_station_id"`  // 下车车站id
	ExpireDate     time.Time `json:"expire_date"`
	State          int       `json:"state"` // 状态，0正在候补，1候补成功，2为候补失败
}

const createdBy = "candidate-server"

// AddCandidate 添加一个订单到数据库
func AddCandidate(can *Candidate) error {
	can.CreatedBy = createdBy
	db := database.Client()
	if err := db.Create(can).Error; err != nil {
		return err
	}
	return nil
}

// AddCandidates 添加多个订单
func AddCandidates(cans []Candidate) error {
	for i := range cans {
		cans[i].CreatedBy = createdBy
	}
	if err := database.Client().Create(&cans).Error; err != nil {
		return err
	}
	return nil
}

// GetCandidates 通过条件获取订单, 条件的key是candidate struct成员
func GetCandidates(conditions map[string]interface{}) ([]*Candidate, error) {
	var res []*Candidate
	if err := database.Client().Where(conditions).Find(&res).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

// GetCandidates 通过用户id查找候补订单
func GetCandidatesByUserID(userID int) ([]*Candidate, error) {
	return GetCandidates(map[string]interface{}{"user_id": userID})
}

// GetCandidatesByTrainID 通过车次id查找
func GetCandidatesByTrainID(trainID int) ([]*Candidate, error) {
	return GetCandidates(map[string]interface{}{"train_id": trainID})
}
