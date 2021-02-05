// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package model

import (
	"pay/model"
	"time"
)

type Candidate struct {
	model.Model
	Date            time.Time `json:"date"`             // 候补时间
	TrainID         uint      `json:"train_id"`         // 车次id
	OrderID         uint      `json:"order_id"`         // 订单id
	UserID          uint      `json:"user_id"`          // 用户id
	PassengerNumber int       `json:"passenger_number"` // 候补人数
	PassengerID     int       `json:"passenger_id"`     // 乘客id
	State           int       `json:"state"`            // 状态，0正在候补，1候补成功，2为候补失败
}
