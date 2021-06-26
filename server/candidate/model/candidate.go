// Package model
// @Author LiuYong
// @Created at 2021-02-04
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
	SeatTypeID     int       `json:"seat_type_id"` // 0商务座, 1一等座, 2二等座
	TicketID       uint      `json:"ticket_id"`
	State          int       `json:"state"` // 状态，0正在候补，1候补成功但是未兑现, 2候补成功，3为候补失败
}

const (
	BusinessSeat = iota
	FirstSeat
	SecondSeat
)

const (
	CandidateIng = iota
	CandidateNotCash
	CandidateSuccess
	CandidateFail
)

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

// GetCandidatesByUserID 通过用户id查找候补订单
func GetCandidatesByUserID(userID int) ([]*Candidate, error) {
	return GetCandidates(map[string]interface{}{"user_id": userID})
}

// GetCandidatesByTrainID 通过车次id查找
func GetCandidatesByTrainID(trainID int) ([]*Candidate, error) {
	return GetCandidates(map[string]interface{}{"train_id": trainID})
}

func GetCandidatesByOrderID(orderID string) ([]*Candidate, error) {
	return GetCandidates(map[string]interface{}{"order_id": orderID})
}

// GetCandidatesOrderIDs 获取正在候补状态的订单, 已创建订单的时间排序
func GetCandidatesOrderIDs() []string {
	selectRes := make([]struct {
		OrderID   string
		CreatedAt time.Time
	}, 0)
	database.Client().Raw("select distinct order_id,created_at from candidates where state = 0 order by created_at asc").Scan(&selectRes)
	res := make([]string, len(selectRes))
	for i := range selectRes {
		res[i] = selectRes[i].OrderID
	}
	return res
}

// GetCandidateTrainIDs 获取需要候补的车次id
func GetCandidateTrainIDs() []uint {
	var trainIDs []uint
	database.Client().Raw("select distinct train_id from candidates").Scan(&trainIDs)
	return trainIDs
}

// UpdateCandidatesState 更新订单中的状态
func UpdateCandidatesState(orderID string, state int) error {
	return database.Client().Table("candidates").Where("order_id = ?", orderID).Updates(map[string]interface{}{"state": state}).Error
}

// UpdateCandidate 更新候补订单
func UpdateCandidate(can *Candidate) error {
	return database.Client().Save(can).Error
}

func GetTrainNumber(trainID uint) string {
	res := ""
	database.Client().Raw("select number from trains where id = ?", trainID).Scan(&res)
	return res
}

func GetStationName(stationID uint) string {
	res := ""
	database.Client().Raw("select name from trains where id = ?", stationID).Scan(&res)
	return res
}
