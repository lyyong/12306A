// @Author KongLingWen
// @Created at 2021/1/29
// @Modified at 2021/1/29

package main

import (
	"common/tools/logging"
	"net"
	"os"
	"os/signal"
	"syscall"
	"ticketPool/rpc"
)

func main() {

	logging.Info("TicketPool Service....")

	/* 初始化 rpc2 (注册rpc服务）*/
	logging.Info("register rpc server")
	rpcServer := rpc.InitRPCServer()

	logging.Info("Listen")
	rpcListen, err := net.Listen("tcp", "127.0.0.1:8002")
	if err != nil {
		logging.Error("listen fail:", err)
		return
	}

	go func(){
		if err := rpcServer.Serve(rpcListen); err != nil {
			logging.Fatal("rpc2 server: ", err)
			return
		}
	}()



	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGKILL)
	<-quit

	logging.Info("TicketPool Server Closed")

}
