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

const (
	ORDER_NOT_FINISH = iota // state = 0 订单未支付完成
	ORDER_FINISH            // state = 1 订单已经完成
	ORDER_CANCEL            // state = 2 订单已经取消
	ORDER_CHANGE            // state = 3 订单已经改签
	ORDER_REFUND            // state = 4 订单已经退款
)

// AddOrder 数据库写入一个订单信息
func AddOrder(o *Order) error {
	o.CreatedBy = "pay-server"
	if err := database.Client().Create(o).Error; err != nil {
		return err
	}
	return nil
}

// GetOrders 通过不同的条件拆查询订单
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

// GetOrdersByUserID 通过用户id查询订单信息
func GetOrdersByUserID(userID int) ([]*Order, error) {
	return GetOrders(map[string]interface{}{"user_id": userID})
}

// GetOrdersByAffairID 通过业务id来查询订单
func GetOrdersByAffairID(affairID string) ([]*Order, error) {
	return GetOrders(map[string]interface{}{"affair_id": affairID})
}

// GetOrderByID 通过订单id查询订单信息
func GetOrderByID(id int) (*Order, error) {
	var order Order
	if err := database.Client().Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetOrderByOutsideID 通过订单的外部id查询订单信息
func GetOrderByOutsideID(outsideID string) (*Order, error) {
	res, err := GetOrders(map[string]interface{}{"outside_id": outsideID})
	if err != nil {
		return nil, err
	}
	if res != nil && len(res) > 0 {
		return res[0], nil
	}
	return nil, nil
}

// UpdateOrder 更新订单信息
func UpdateOrder(order *Order) error {
	if err := database.Client().Model(order).Updates(*order).Error; err != nil {
		return err
	}
	return nil
}
