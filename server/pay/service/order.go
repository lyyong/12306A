// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package service

import (
	"common/tools/logging"
	"errors"
	"pay/model"
)

type OrderService struct {
}

// CreateOrder 创建一个订单
func (s OrderService) CreateOrder(userID uint, money, affairID, createdBy string) (string, error) {
	order := model.Order{
		Model:    model.Model{CreatedBy: createdBy},
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

func (s OrderService) UpdateOrderState(outsideID string, state int) error {
	order, err := model.GetOrderByOutsideID(outsideID)
	if err != nil {
		return err
	}
	if order == nil {
		err = errors.New("订单id错误")
		logging.Error(err)
		return err
	}
	order.State = state
	return nil
}

func (s OrderService) UpdateOrderStateWithRelative(outsideID string, state int, relativeID string) error {
	order, err := model.GetOrderByOutsideID(outsideID)
	if err != nil {
		logging.Error(err)
		return err
	}
	if order == nil {
		err = errors.New("订单id错误")
		logging.Error(err)
		return err
	}
	rorder, err := model.GetOrderByOutsideID(relativeID)
	if err != nil {
		logging.Error(err)
		return err
	}
	if order == nil {
		err = errors.New("关联订单id错误")
		logging.Error(err)
		return err
	}
	order.RelativeOrder = rorder.ID
	order.State = state
	model.UpdateOrder(order)
	if err != nil {
		logging.Error(err)
		return err
	}
	return nil
}

func (s OrderService) GetOrdersByUserID(userID int) []*model.Order {
	orders, err := model.GetOrdersByUserID(userID)
	if err != nil {
		logging.Error(err)
		return nil
	}
	if orders == nil || len(orders) == 0 {
		return nil
	}
	return orders
}
