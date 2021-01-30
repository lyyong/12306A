// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package rpc

import (
	"context"
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

func (o OrderRPCService) Create(ctx context.Context, info *orderRPCpb.CreateInfo) (*orderRPCpb.Error, error) {
	return &orderRPCpb.Error{Content: "hello Create"}, nil
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
