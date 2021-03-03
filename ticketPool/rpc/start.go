// @Author LiuYong
// @Created at 2021-02-12
// @Modified at 2021-02-12
package rpc

import (
	"common/router_tracer"
	"common/rpc_manage"
	"common/server_find"
	"common/tools/logging"
	"fmt"
	"net"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"ticketPool/controller"
	"ticketPool/tools/setting"
)

func Setup() {
	// 载入配置文件
	setting.Setup()
	// 配置服务发现
	server_find.Register(setting.Server.Name, setting.Server.Host, fmt.Sprintf("%d", setting.Server.HttpPort), setting.Consul.ServiceID, setting.Consul.Address, setting.Consul.Interval, setting.Consul.TTL)
	// 链路追踪
	router_tracer.SetupByHttp(setting.Zipkin.ServiceID, setting.Server.Host, fmt.Sprintf("%d", setting.Server.HttpPort), setting.Zipkin.HttpEndpoint)
	// 开启rpc
	rpcServer := rpc_manage.NewGRPCServer()
	ticketPoolRPC.RegisterTicketPoolServiceServer(rpcServer, &controller.TicketPoolRPCService{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", setting.Server.HttpPort))
	if err != nil {
		logging.Error(err)
	}

	//go func() {
	//	if err := rpcServer.Serve(lis); err != nil {
	//		logging.Fatal("rpc server: ", err)
	//		return
	//	}
	//}()
	rpcServer.Serve(lis)
	defer lis.Close()
}
