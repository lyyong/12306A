// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package service

import (
	"common/tools/logging"
	"errors"
	"pay/model"
	"pay/service/cache"
	cache2 "pay/tools/cache"
)

type OrderService struct {
}

const orderExpTime = 1800

// CreateOrder 创建一个订单
func (s OrderService) CreateOrder(userID uint, money int, affairID, createdBy string) (string, error) {
	order := model.Order{
		Model:    model.Model{CreatedBy: createdBy},
		UserID:   userID,
		Money:    money,
		AffairID: affairID,
	}
	order.OutsideID = order.AffairID
	orderCache := cache.OrderCache{
		UserID:    userID,
		OutsideID: order.OutsideID,
		AffairID:  order.AffairID,
	}
	// 添加到redis设置30分钟期限
	cache2.Set(orderCache.GetNoFinishOrderKey(), order, orderExpTime)

	// 存入数据库
	// defer func() {
	// 	if err := model.AddOrder(&order); err != nil {
	// 		logging.Error(err)
	// 	}
	// }()
	return order.OutsideID, nil
}

func (s OrderService) UpdateOrderState(outsideID string, state int) error {
	order := getOrderByOutsideID(outsideID)
	var err error = nil
	if order == nil {
		err = errors.New("订单id错误")
		logging.Error(err)
		return err
	}
	saveOrderWithStateChange(order, state)
	return nil
}

func (s OrderService) UpdateOrderStateWithRelative(outsideID string, state int, relativeID string) error {
	order := getOrderByOutsideID(outsideID)
	if order == nil {
		err := errors.New("订单id错误")
		logging.Error(err)
		return err
	}
	rorder := getOrderByOutsideID(relativeID)
	if rorder == nil {
		err := errors.New("关联订单id错误")
		logging.Error(err)
		return err
	}
	order.RelativeOrder = rorder.ID
	saveOrderWithStateChange(order, state)
	return nil
}

// GetOrdersByUserID 得到已完成的订单
func (s OrderService) GetOrdersByUserID(userID uint) []*model.Order {
	orderCache := cache.OrderCache{
		UserID: userID,
	}
	if cache2.Exists(orderCache.GetOrdersKey()) {
		orders := make([]*model.Order, 0)
		err := cache2.Get(orderCache.GetOrdersKey(), &orders)
		if err == nil {
			return orders
		}
	}
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

func (s OrderService) GetOrdersByUserIDAndUnfinish(userID uint) *model.Order {
	orderCache := cache.OrderCache{
		UserID: userID,
	}
	if cache2.Exists(orderCache.GetNoFinishOrderKey()) {
		var order model.Order
		err := cache2.Get(orderCache.GetNoFinishOrderKey(), &order)
		if err == nil {
			return &order
		}
	}
	return nil
}

// getOrderByOutsideID 通过redis或者mysql获得order
func getOrderByOutsideID(outsideID string) *model.Order {
	var order *model.Order
	var err error = nil
	orderCache := cache.OrderCache{
		OutsideID: outsideID,
	}
	if cache2.Exists(orderCache.GetOrderKey()) {
		order = new(model.Order)
		err = cache2.Get(orderCache.GetOrderKey(), order)
		if err != nil {
			logging.Error(err)
			order, err = model.GetOrderByOutsideID(outsideID)
			if err != nil {
				return nil
			}
		}
	}
	order, err = model.GetOrderByOutsideID(outsideID)
	if err != nil {
		return nil
	}
	return order
}

// saveOrderWithStateChange 保存状态变更的order
func saveOrderWithStateChange(order *model.Order, state int) {
	orderCache := cache.OrderCache{
		UserID:    order.UserID,
		OutsideID: order.OutsideID,
		AffairID:  order.AffairID,
	}
	if order.State == 0 {
		// 未完成的订单需要存入数据库
		defer func() {
			model.AddOrder(order)
		}()
		order.State = state
		cache2.Set(orderCache.GetOrderKey(), order, orderExpTime)
		return
	}
	order.State = state
	cache2.Set(orderCache.GetOrderKey(), order, orderExpTime)
	defer func() {
		model.UpdateOrder(order)
	}()
}
