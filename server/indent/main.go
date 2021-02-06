// @Author KongLingWen
// @Created at 2021/1/29
// @Modified at 2021/1/29

package main

import (
	"common/tools/logging"
	"context"
	"indent/routers"
	"indent/rpc"
	"indent/utils/setting"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	logging.Info("Indent Service....")

	/* 初始化 rpc (注册rpc服务）*/
	logging.Info("register rpc server")
	rpcServer := rpc.InitRPCServer()

	rpcListen, err := net.Listen("tcp", setting.Server.RpcAddr)
	if err != nil {
		logging.Error("listen fail: ", err)
		return
	}

	go func(){
		if err := rpcServer.Serve(rpcListen); err != nil {
			logging.Fatal("rpc server: ", err)
			return
		}
	}()

	/* 初始化 router */
	initRouter := routers.InitRouter()
	server := &http.Server{
		Addr:              setting.Server.HttpAddr,
		Handler:           initRouter,
	}
	go func(){
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			logging.Fatal("listen fail: ", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGKILL)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		logging.Fatal("Server Shutdown:", err)
	}
	logging.Info("Indent Server Closed")

}
