// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package service

import (
	"common/tools/logging"
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"pay/model"
	"pay/service/cache"
	cache2 "pay/tools/cache"
	"pay/tools/setting"
	"rpc/pay/client/orderRPCClient"
	"strconv"
	"sync"
	"time"
)

type orderService struct {
}

var (
	kwToOrder     *kafka.Writer
	onceForKO     sync.Once
	topicRefundOK = "RefundOK"
)

const (
	orderExpTime = 1800
)

func NewOrderService() *orderService {
	kwToOrder = getKWToOrder()
	return &orderService{}
}

// CreateOrder 创建一个订单
func (s orderService) CreateOrder(userID uint, money int, affairID, createdBy string) (string, error) {
	order := model.Order{
		Model:    model.Model{CreatedBy: createdBy},
		UserID:   userID,
		Money:    money,
		AffairID: affairID,
	}
	order.OutsideID = createdBy + time.Now().Format("2006-01-02-15:04:05.000-") + strconv.Itoa(int(userID))
	orderCache := cache.OrderCache{
		UserID:    userID,
		OutsideID: order.OutsideID,
	}
	// 添加到redis设置30分钟期限
	cache2.Set(orderCache.GetUnpayOrderKey(), order, orderExpTime)

	// 存入数据库
	// defer func() {
	// 	if err := model.AddOrder(&order); err != nil {
	// 		logging.Error(err)
	// 	}
	// }()
	return order.OutsideID, nil
}

func (s orderService) UpdateOrderState(outsideID string, state int) error {
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

func (s orderService) UpdateOrderStateWithRelative(outsideID string, state int, relativeID string) error {
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

// GetOrdersByUserIDAndFinished 得到用户的历史订单
func (s orderService) GetOrdersByUserIDAndFinished(userID uint) []*model.Order {
	return s.GetOrdersByUserIDAndState(userID, model.ORDER_FINISHED)
}

// GetOrdersByUserIDAndUnpay 获取用户未支付的订单
func (s orderService) GetOrdersByUserIDAndUnpay(userID uint) *model.Order {
	orderCache := cache.OrderCache{
		UserID: userID,
	}
	if cache2.Exists(orderCache.GetUnpayOrderKey()) {
		var order model.Order
		err := cache2.Get(orderCache.GetUnpayOrderKey(), &order)
		if err == nil {
			return &order
		}
	}
	return nil
}

// GetOrdersByUserIDAndUnfinished 获取没有出行的订单
func (s orderService) GetOrdersByUserIDAndUnfinished(userID uint) []*model.Order {
	return s.GetOrdersByUserIDAndState(userID, model.ORDER_UNFINISHED)
}

// GetOrdersByUserIDAndState 通过用户名和状态获取订单
func (s orderService) GetOrdersByUserIDAndState(userID uint, state int) []*model.Order {
	orderCache := cache.OrderCache{
		UserID: userID,
	}
	if cache2.Exists(orderCache.GetOrdersKey()) {
		orders := make([]*model.Order, 0)
		err := cache2.Get(orderCache.GetOrdersKey(), &orders)
		if err == nil {
			// 保留已完成的订单
			for i := 0; i < len(orders); i++ {
				if orders[i].State == state {
					continue
				}
				orders = append(orders[:i], orders[i+1:]...)
				i--
			}
			return orders
		}
	}
	orders, err := model.GetOrders(map[string]interface{}{"user_id": userID, "state": state})
	if err != nil {
		logging.Error(err)
		return nil
	}
	if orders == nil || len(orders) == 0 {
		return nil
	}
	return orders
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

// Refund 订单退款
func (s orderService) Refund(userID uint, outsideID string, fullMoney bool, money int) error {
	orderCache := cache.OrderCache{
		UserID:    userID,
		OutsideID: outsideID,
	}
	orders := make([]*model.Order, 0)
	if cache2.Exists(orderCache.GetOrdersKey()) {
		cache2.Get(orderCache.GetOrdersKey(), &orders)
	} else {
		orders, _ = model.GetOrdersByUserID(userID)
	}
	for _, order := range orders {
		if order.OutsideID == outsideID {
			order.State = model.ORDER_REFUND
			// TODO 真实退款
			cache2.Delete(orderCache.GetOrdersKey())
			cache2.Set(orderCache.GetOrdersKey(), orders, expTime)
			model.UpdateOrder(order)
			// 消息队列发送取消订单完成
			refundOKInfo := orderRPCClient.RefundOKInfo{OutsideID: outsideID}
			msgV, err := json.Marshal(&refundOKInfo)
			if err != nil {
				logging.Error(err)
				return err
			}
			err = kwToPay.WriteMessages(context.TODO(), kafka.Message{
				Value: msgV,
			})
			if err != nil {
				logging.Error(err)
				return err
			}
			return nil
		}
	}
	return errors.New("无此订单")
}

func (s *orderService) CancelUnpayOrder(userID uint) {
	orderCache := cache.OrderCache{
		UserID: userID,
	}
	cache2.Delete(orderCache.GetUnpayOrderKey())
}

func getKWToOrder() *kafka.Writer {
	onceForKO.Do(func() {
		if kwToOrder == nil {
			kwToOrder = &kafka.Writer{
				Addr:         kafka.TCP(setting.Kafka.Host),
				Topic:        topicRefundOK,
				Async:        false,                 // 非异步执行
				BatchTimeout: 50 * time.Millisecond, // 消息在发送缓存中的等待时间, 设置小点速度快但是占用cpu
				WriteTimeout: 10 * time.Second,
			}
		}
	})
	return kwToOrder
}
