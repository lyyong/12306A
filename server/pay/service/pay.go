// @Author LiuYong
// @Created at 2021-02-08
// @Modified at 2021-02-08
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
	"sync"
	"time"
)

type payService struct {
}

const (
	orderInfo = "asdiuyUYGFYGV7567hgvfhjv"
	expTime   = 1800
)

var (
	kwToPay    *kafka.Writer
	onceForKP  sync.Once
	topicPayOK = "PayOK"
)

func NewPayService() *payService {
	kwToPay = getKWToPay()
	return &payService{}
}

// WantPay 用户准备支付
// userID 用户的ID
// orderOutsideID 订单外部IDu
// 返回OrderInfo
func (s payService) WantPay(userID uint, orderOutsideID string) string {
	orderCache := cache.OrderCache{
		UserID:    userID,
		OutsideID: orderOutsideID,
	}
	if cache2.Exists(orderCache.GetUnpayOrderKey()) {
		var order model.Order
		cache2.Get(orderCache.GetUnpayOrderKey(), &order)
		if order.OutsideID != orderOutsideID {
			logging.Error("outsideID不符合", order.OutsideID, "!=", orderOutsideID)
			return ""
		}
		order.AlipayOrderInfo = orderInfo
		cache2.Delete(orderCache.GetUnpayOrderKey())
		cache2.Set(orderCache.GetUnpayOrderKey(), &order, expTime)
		return orderInfo
	}
	return ""
}

func (s payService) PayOK(userID uint, orderInfo, orderOutsideID string) error {
	orderCache := cache.OrderCache{
		UserID:    userID,
		OutsideID: orderOutsideID,
	}
	if !cache2.Exists(orderCache.GetUnpayOrderKey()) {
		return errors.New("订单信息出错")
	}
	var order model.Order
	cache2.Get(orderCache.GetUnpayOrderKey(), &order)
	order.State = model.ORDER_UNFINISHED
	// TODO 下面的代码也许可以放到defer中运行
	payOKInfo := orderRPCClient.PayOKOrderInfo{
		UserID:    order.UserID,
		OutsideID: order.OutsideID,
		AffairID:  order.AffairID,
		Money:     order.Money,
		State:     order.State,
	}
	// 发送支付完成消息到消息队列
	msgV, err := json.Marshal(&payOKInfo)
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
	// 更新到数据库
	model.AddOrder(&order)
	// 删除cache中的未完成订单
	cache2.Delete(orderCache.GetUnpayOrderKey())
	// 更新cache
	if !cache2.Exists(orderCache.GetOrdersKey()) {
		return nil
	}
	orders := make([]*model.Order, 0)
	cache2.Get(orderCache.GetOrdersKey(), &order)
	orders = append(orders, &order)
	cache2.Set(orderCache.GetOrdersKey(), &orders, expTime)
	return nil
}

func getKWToPay() *kafka.Writer {
	onceForKP.Do(func() {
		if kwToPay == nil {
			kwToPay = &kafka.Writer{
				Addr:         kafka.TCP(setting.Kafka.Host),
				Topic:        topicPayOK,
				Async:        false,                 // 非异步执行
				BatchTimeout: 50 * time.Millisecond, // 消息在发送缓存中的等待时间, 设置小点速度快但是占用cpu
				WriteTimeout: 10 * time.Second,
			}
		}
	})
	return kwToPay
}
