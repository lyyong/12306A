// @Author LiuYong
// @Created at 2021-01-27
// @Modified at 2021-01-27
package router

import (
	"common/rpcServer"
	"pay/controller/rpc"
	"pay/tools/setting"
	"rpc/pay/proto/orderRPCpb"
)

func InitRPCService() (*rpcServer.RpcServer, error) {
	server, err := rpcServer.NewRpcServerWithServerFindAndHttpTracer(setting.Server.Name,
		setting.Server.Host, setting.Server.HttpPort, setting.Consul.ServiceID, setting.Consul.Address, setting.Consul.Interval, setting.Consul.TTL, setting.Zipkin.HttpEndpoint)
	if err != nil {
		return nil, err
	}
	orderRPCpb.RegisterOrderRPCServiceServer(server.Server(), &rpc.OrderRPCService{})
	return server, nil
}
