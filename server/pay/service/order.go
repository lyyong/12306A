// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package service

import "pay/model"

type OrderService struct {
}

func (s OrderService) AddOrder(userID uint, money string, affairID string) (string, error) {
	order := model.Order{
		UserID:   userID,
		Money:    money,
		AffairID: affairID,
	}
	order.OutsideID = order.AffairID
	if err := model.AddOrder(&order); err != nil {
		return "", err
	}
	return order.OutsideID, nil
}
