// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCImp

import (
	"errors"
	"interface/pay/orderInterfaces"
	"rpc/pay/client/orderRPCClient"
	"rpc/pay/proto/orderRPCpb"
)

type OrderRPCImp struct {
	client *orderRPCClient.OrderRPCClient
}

func NewOrderRPCImp() (*OrderRPCImp, error) {
	cli, err := orderRPCClient.NewClient()
	if err != nil {
		return nil, err
	}
	return &OrderRPCImp{client: cli}, nil
}

// Create 创建订单
func (o *OrderRPCImp) Create(info *orderInterfaces.CreateInfo) error {
	res, err := o.client.Create(&orderRPCpb.CreateInfo{
		UserID:         info.UserID,
		Money:          info.Money,
		AffairID:       info.AffairID,
		ExpireDuration: info.ExpireDuration,
		OrderOutsideID: info.OrderOutsideID,
	})
	if err != nil {
		return err
	}
	if res != nil && res.Content != "" {
		return errors.New(res.Content)
	}
	return nil
}

// Read 获取订单
// userID 为用户的逐渐
func (o *OrderRPCImp) Read(userID int64) (*orderInterfaces.Info, error) {
	res, err := o.client.Read(&orderRPCpb.SearchInfo{UserID: userID})
	if err != nil {
		return nil, err
	}
	return &orderInterfaces.Info{
		UserID:         res.UserID,
		Money:          res.Money,
		AffairID:       res.AffairID,
		ExpireDuration: res.ExpireDuration,
		OrderOutsideID: res.OrderOutsideID,
		State:          res.State,
	}, nil
}

func (o *OrderRPCImp) UpdateState(orderOutsideID string, state int32) error {
	res, err := o.client.UpdateState(&orderRPCpb.UpdateStateInfo{
		OutsideID: orderOutsideID,
		State:     state,
	})
	if err != nil {
		return err
	}
	if res != nil && res.Content != "" {
		return errors.New(res.Content)
	}
	return nil
}

func (o *OrderRPCImp) UpdateStateWithRelativeOrder(orderOutsideID string, state int32, relativeOutsideID string) error {
	res, err := o.client.UpdateStateWithRelativeOrder(&orderRPCpb.UpdateStateWithRInfo{
		OutsideID:  orderOutsideID,
		State:      state,
		ROutsideID: relativeOutsideID,
	})
	if err != nil {
		return err
	}
	if res != nil && res.Content != "" {
		return errors.New(res.Content)
	}
	return nil
}
