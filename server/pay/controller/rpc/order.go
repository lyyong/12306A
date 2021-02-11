// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package rpc

import (
	"context"
	"pay/service"
	"rpc/pay/proto/orderRPCpb"
)

type OrderRPCService struct {
}

func (o OrderRPCService) UpdateState(ctx context.Context, info *orderRPCpb.UpdateStateInfo) (*orderRPCpb.Error, error) {
	return &orderRPCpb.Error{Content: "hello UpdateState"}, nil
}

func (o OrderRPCService) UpdateStateWithRelativeOrder(ctx context.Context, info *orderRPCpb.UpdateStateWithRInfo) (*orderRPCpb.Error, error) {
	return &orderRPCpb.Error{Content: "hello UpdateStateWithRelativeOrder"}, nil
}

func (o OrderRPCService) Create(ctx context.Context, info *orderRPCpb.CreateInfo) (*orderRPCpb.CreateRes, error) {
	orderService := &service.OrderService{}
	outsideID, err := orderService.AddOrder(uint(info.UserID), info.Money, info.AffairID)
	if err != nil {
		return nil, err
	}
	return &orderRPCpb.CreateRes{OrderOutsideID: outsideID}, nil
}

func (o OrderRPCService) Read(ctx context.Context, info *orderRPCpb.SearchInfo) (*orderRPCpb.Info, error) {
	return &orderRPCpb.Info{
		UserID:         0,
		Money:          "30",
		AffairID:       "123zsd",
		ExpireDuration: 0,
		OrderOutsideID: "",
		State:          0,
	}, nil
}
