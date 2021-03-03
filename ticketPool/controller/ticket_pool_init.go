/*
* @Author: 余添能
* @Date:   2021/3/3 10:11 下午
 */
package controller

import (
	"common/rpc_manage"
	"common/tools/logging"
	"google.golang.org/grpc"
	"rpc/ticketPool/proto/ticketPoolRPC"
)

func InitRPCServer() *grpc.Server {
	server := rpc_manage.NewGRPCServer()
	logging.Info("Register TicketPoolService Server")
	ticketPoolRPC.RegisterTicketPoolServiceServer(server, new(TicketPoolRPCService))
	return server
}
