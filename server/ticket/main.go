// @Author KongLingWen
// @Created at 2021/1/29
// @Modified at 2021/1/29

package main

import (
	"common/router_tracer"
	"common/server_find"
	"common/tools/logging"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"ticket/routers"
	"ticket/rpc"
	"ticket/utils/database"
	"ticket/utils/redispool"
	"ticket/utils/setting"
)

func main() {
	logging.Info("Ticket Service Start....")

	server_find.Register(setting.Server.Name,
		setting.Server.Host, strconv.Itoa(setting.Server.HttpPort), setting.Consul.ServiceID, setting.Consul.Address, setting.Consul.Interval, setting.Consul.TTL)
	// 链路追踪
	err := router_tracer.SetupByHttp(setting.Server.Name,
		setting.Server.Host, strconv.Itoa(setting.Server.HttpPort), setting.Zipkin.HttpEndpoint)

	/* 初始化 rpc (注册rpc服务）*/
	logging.Info("register rpc server")
	rpcServer := rpc.InitRPCServer()
	logging.Info("Listen:",setting.Server.RPCAddr)
	rpcListen, err := net.Listen("tcp", setting.Server.RPCAddr)
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

	logging.Info("Close database connection pool...")
	sqlDB, err := database.DB.DB()
	if err == nil {
		sqlDB.Close()
	}

	logging.Info("Close redis connection pool...")
	redispool.RedisPool.Close()

	if err := server.Shutdown(context.Background()); err != nil {
		logging.Fatal("Ticket Server Shutdown error:", err)
	}
	logging.Info("Ticket Server Closed")

}
