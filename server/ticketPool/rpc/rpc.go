package rpc

import (
	"common/rpc_manage"
	"common/tools/logging"
	"google.golang.org/grpc"
	"rpc/ticketPool/proto/ticketPoolRPC"
)

func InitRPCServer() *grpc.Server {
	server := rpc_manage.NewGRPCServer()
	logging.Info("Register TicketPoolService Server")
	ticketPoolRPC.RegisterTicketPoolServiceServer(server, &TicketPoolServer{})
	return server
}
