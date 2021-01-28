// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package router

import (
	"google.golang.org/grpc"
	"pay/controller/rpc"
	"pay/tools/logging"
	"rpc/pay/proto/orderRPCpb"
)

func InitRPCService() *grpc.Server {
	server := grpc.NewServer()
	logging.Info("注册OrderRPCService作为RPC服务器")
	orderRPCpb.RegisterOrderRPCServiceServer(server, &rpc.OrderRPCService{})
	return server
}
