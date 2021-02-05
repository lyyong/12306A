// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package model

import "pay/tools/database"

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

func DeleteOrder(condition map[string]interface{}) ([]*Order, error) {
	var res []*Order
	if err := database.Client().Where(condition).Find(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
