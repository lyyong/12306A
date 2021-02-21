// @Author LiuYong
// @Created at 2021-02-08
// @Modified at 2021-02-08
package service

import "pay/model"

type PayService struct {
}

// WantPay 用户准备支付
// userID 用户的ID
// orderOutsideID 订单外部IDu
// 返回OrderInfo
func (s PayService) WantPay(userID uint, orderOutsideID string) string {
	orders, err := model.GetOrdersByUserID(userID)
	if err != nil {
		return ""
	}
	var order *model.Order = nil
	for _, o := range orders {
		if o.OutsideID == orderOutsideID {
			order = o
		}
	}
	if order == nil {
		return ""
	}

	orderInfo := "asdiuyUYGFYGV7567hgvfhjv"

	order.AlipayOrderInfo = orderInfo
	model.UpdateOrder(order)
	return orderInfo
}
