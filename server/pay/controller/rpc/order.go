// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package rpc

import (
	"context"
	"errors"
	"pay/service"
	"rpc/pay/proto/orderRPCpb"
)

type OrderRPCService struct {
}

func (o OrderRPCService) Refund(ctx context.Context, request *orderRPCpb.RefundRequest) (*orderRPCpb.Respond, error) {
	orderService := service.NewOrderService()
	return nil, orderService.Refund(uint(request.UserID), request.OutsideID, request.FullMoney, int(request.Money))
}

func (o OrderRPCService) GetNoFinishOrder(ctx context.Context, condition *orderRPCpb.SearchCondition) (*orderRPCpb.OrderInfo, error) {
	orderService := service.NewOrderService()
	order := orderService.GetOrdersByUserIDAndUnfinish(uint(condition.UserID))
	if order == nil {
		return nil, nil
	}
	return &orderRPCpb.OrderInfo{
		UserID:         uint64(order.UserID),
		Money:          int64(order.Money),
		AffairID:       order.AffairID,
		ExpireDuration: int32(order.ExpireDuration),
		OrderOutsideID: order.OutsideID,
		State:          int32(order.State),
	}, nil
}

// UpdateState RPC更新订单状态
func (o OrderRPCService) UpdateState(ctx context.Context, info *orderRPCpb.UpdateStateRequest) (*orderRPCpb.Respond, error) {
	orderService := service.NewOrderService()
	err := orderService.UpdateOrderState(info.OutsideID, int(info.State))
	if err != nil {
		return &orderRPCpb.Respond{Content: err.Error()}, nil
	}
	return nil, nil
}

// UpdateStateWithRelativeOrder RPC更新订单状态添加相关订单
func (o OrderRPCService) UpdateStateWithRelativeOrder(ctx context.Context, info *orderRPCpb.UpdateStateWithRRequest) (*orderRPCpb.Respond, error) {
	orderService := service.NewOrderService()
	err := orderService.UpdateOrderStateWithRelative(info.OutsideID, int(info.State), info.ROutsideID)
	if err != nil {
		return &orderRPCpb.Respond{Content: err.Error()}, nil
	}
	return &orderRPCpb.Respond{Content: "hello UpdateStateWithRelativeOrder"}, nil
}

// Create RPC创建订单
func (o OrderRPCService) Create(ctx context.Context, info *orderRPCpb.CreateRequest) (*orderRPCpb.CreateRespond, error) {
	orderService := service.NewOrderService()
	// 判断该用户时是否有未完成的订单
	order := orderService.GetOrdersByUserIDAndUnfinish(uint(info.UserID))
	if order != nil {
		return nil, errors.New("客户存在未完成的订单")
	}
	outsideID, err := orderService.CreateOrder(uint(info.UserID), int(info.Money), info.AffairID, info.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &orderRPCpb.CreateRespond{OrderOutsideID: outsideID}, nil
}

// Read 获取用户的相关订单
func (o OrderRPCService) Read(ctx context.Context, info *orderRPCpb.SearchCondition) (*orderRPCpb.ReadRespond, error) {
	orderService := service.NewOrderService()
	orders := orderService.GetOrdersByUserID(uint(info.UserID))
	var readRespond orderRPCpb.ReadRespond
	for _, order := range orders {
		readRespond.Infos = append(readRespond.Infos, &orderRPCpb.OrderInfo{
			UserID:         uint64(order.UserID),
			Money:          int64(order.Money),
			AffairID:       order.AffairID,
			ExpireDuration: int32(order.ExpireDuration),
			OrderOutsideID: order.OutsideID,
			State:          int32(order.State),
		})
	}
	return &readRespond, nil
}
