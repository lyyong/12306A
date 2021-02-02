// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package router

import (
	"common/rpc_manage"
	"google.golang.org/grpc"
	"pay/controller/rpc"
	"rpc/pay/proto/orderRPCpb"
)

func InitRPCService() *grpc.Server {
	server := rpc_manage.NewGRPCServer()
	orderRPCpb.RegisterOrderRPCServiceServer(server, &rpc.OrderRPCService{})
	return server
}
