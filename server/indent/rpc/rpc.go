package rpc

import (
	"common/rpc_manage"
	"common/tools/logging"
	"google.golang.org/grpc"
	"rpc/indent/proto/indentRPC"
)

func InitRPCServer() *grpc.Server {
	//server := grpc.NewServer()
	//logging.Info("Register IndentService Server")
	//indentRPC.RegisterIndentServiceServer(server, &IndentServer{})
	server := rpc_manage.NewGRPCServer()
	logging.Info("Register IndentService Server")
	indentRPC.RegisterIndentServiceServer(server, &IndentServer{})
	return server
}
