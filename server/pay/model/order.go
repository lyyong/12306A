// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package model

import (
	"github.com/jinzhu/gorm"
	"pay/tools/database"
)

type Order struct {
	Model
	UserID          uint   `json:"user_id"`
	AlipayOrderInfo string `json:"alipay_order_info"`
	Money           string `json:"money"`
	AffairID        string `json:"affair_id"`
	ExpireDuration  int    `json:"expire_duration"`
	OutsideID       string `json:"outside_id"`
	RelativeOrder   uint   `json:"relative_order"`
	State           int    `json:"state"`
}

func AddOrder(o *Order) error {
	o.CreatedBy = "pay-server"
	if err := database.Client().Create(o).Error; err != nil {
		return err
	}
	return nil
}

func GetOrders(condition map[string]interface{}) ([]*Order, error) {
	var res []*Order
	if err := database.Client().Where(condition).Find(res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func GetOrdersByUserID(userID int) ([]*Order, error) {
	return GetOrders(map[string]interface{}{"user_id": userID})
}

func GetOrdersByAffairID(affairID string) ([]*Order, error) {
	return GetOrders(map[string]interface{}{"affair_id": affairID})
}

func UpdateOrder(order *Order, info map[string]interface{}) error {
	if err := database.Client().Model(order).Updates(info).Error; err != nil {
		return err
	}
	return nil
}
