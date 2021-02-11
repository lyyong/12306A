package rpc

import (
	"common/rpc_manage"
	"common/tools/logging"
	"google.golang.org/grpc"
	"rpc/ticket/proto/ticketRPC"
)

func InitRPCServer() *grpc.Server {
	server := rpc_manage.NewGRPCServer()
	logging.Info("Register IndentService Server")
	ticketRPC.RegisterTicketServiceServer(server, &TicketServer{})
	return server
}
