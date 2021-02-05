// @Author KongLingWen
// @Created at 2021/1/29
// @Modified at 2021/1/29

package main

import (
	"common/tools/logging"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ticket/routers"
	"ticket/rpc"
)

func main() {

	logging.Info("Ticket Service Start....")

	/* 初始化 rpc (注册rpc服务）*/
	logging.Info("register rpc server")
	rpcServer := rpc.InitRPCServer()

	logging.Info("Listen")
	rpcListen, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		logging.Error("listen fail:", err)
		return
	}

	go func(){
		if err := rpcServer.Serve(rpcListen); err != nil {
			logging.Fatal("rpc server: ", err)
			return
		}
	}()

	/* 初始化 router */
	logging.Info("Register Router")
	initRouter := routers.InitRouter()
	server := &http.Server{
		Addr:              ":8080",
		Handler:           initRouter,
	}
	go func(){
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			logging.Fatal("listen: ", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGKILL)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		logging.Fatal("Ticket Server Shutdown:", err)
	}
	logging.Info("Ticket Server Closed")

}
