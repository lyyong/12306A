// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package rpc

import (
	"context"
	"errors"
	"pay/model"
	"pay/service"
	"rpc/pay/proto/orderRPCpb"
)

type OrderRPCService struct {
}

// UpdateState RPC更新订单状态
func (o OrderRPCService) UpdateState(ctx context.Context, info *orderRPCpb.UpdateStateInfo) (*orderRPCpb.Error, error) {
	orderService := service.OrderService{}
	err := orderService.UpdateOrderState(info.OutsideID, int(info.State))
	if err != nil {
		return &orderRPCpb.Error{Content: err.Error()}, nil
	}
	return nil, nil
}

// UpdateStateWithRelativeOrder RPC更新订单状态添加相关订单
func (o OrderRPCService) UpdateStateWithRelativeOrder(ctx context.Context, info *orderRPCpb.UpdateStateWithRInfo) (*orderRPCpb.Error, error) {
	orderService := service.OrderService{}
	err := orderService.UpdateOrderStateWithRelative(info.OutsideID, int(info.State), info.ROutsideID)
	if err != nil {
		return &orderRPCpb.Error{Content: err.Error()}, nil
	}
	return &orderRPCpb.Error{Content: "hello UpdateStateWithRelativeOrder"}, nil
}

// Create RPC创建订单
func (o OrderRPCService) Create(ctx context.Context, info *orderRPCpb.CreateInfo) (*orderRPCpb.CreateRes, error) {
	orderService := &service.OrderService{}
	// 判断该用户时是否有未完成的订单
	orders := orderService.GetOrdersByUserID(uint(info.UserID))
	for _, order := range orders {
		if order.State == model.ORDER_NOT_FINISH {
			return nil, errors.New("客户存在未完成的订单")
		}
	}
	outsideID, err := orderService.CreateOrder(uint(info.UserID), info.Money, info.AffairID, info.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &orderRPCpb.CreateRes{OrderOutsideID: outsideID}, nil
}

// Read 获取用户的相关订单
func (o OrderRPCService) Read(ctx context.Context, info *orderRPCpb.SearchInfo) (*orderRPCpb.ReadInfo, error) {
	orderService := &service.OrderService{}
	orders := orderService.GetOrdersByUserID(uint(info.UserID))
	var readInfo orderRPCpb.ReadInfo
	for _, order := range orders {
		readInfo.Infos = append(readInfo.Infos, &orderRPCpb.OrderInfo{
			UserID:         uint64(order.UserID),
			Money:          order.Money,
			AffairID:       order.AffairID,
			ExpireDuration: int32(order.ExpireDuration),
			OrderOutsideID: order.OutsideID,
			State:          int32(order.State),
		})
	}
	return &readInfo, nil
}
