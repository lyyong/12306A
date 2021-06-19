// @Author KongLingWen
// @Created at 2021/1/29
// @Modified at 2021/1/29

package main

import (
	"common/router_tracer"
	"common/server_find"
	"common/tools/logging"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"ticketPool/rpc"
	"ticketPool/serialize"
	"ticketPool/ticketpool"
	"ticketPool/utils/database"
	"ticketPool/utils/setting"
)

func init() {
	setting.InitSetting()
	if setting.Server.RunMode == "debug" {
		logging.SetupWithMode(logging.LogDebug)
	} else {
		logging.SetupWithMode(logging.LogRelease)
	}
	database.Setup()
}

func Close() {
	database.Close()
}

func main() {

	logging.Info("TicketPool Service start....")
	server_find.Register(setting.Server.Name,
		setting.Server.Host, strconv.Itoa(setting.Server.RPCPort), setting.Consul.ServiceID, setting.Consul.Address, setting.Consul.Interval, setting.Consul.TTL)
	// 链路追踪
	err := router_tracer.SetupByHttp(setting.Server.Name,
		setting.Server.Host, strconv.Itoa(setting.Server.RPCPort), setting.Zipkin.HttpEndpoint)


	/* 初始化 rpc (注册rpc服务）*/
	logging.Info("register rpc server")
	rpcServer := rpc.InitRPCServer()

	logging.Info("Listen", setting.Server.RPCAddr)
	rpcListen, err := net.Listen("tcp", setting.Server.RPCAddr)
	if err != nil {
		logging.Error("listen fail:", err)
		return
	}

	go func() {
		if err := rpcServer.Serve(rpcListen); err != nil {
			logging.Fatal("rpc server: ", err)
			return
		}
	}()

	/* 初始化票池 */
	logging.Info("Init TicketPool")
	ticketpool.InitTicketPool()
	logging.Info("TicketPool init success")

	/* 开启票池序列化 */
	serialize.Serialize()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGKILL)
	<-quit

	logging.Info("TicketPool Server Closed")
	defer func() {
		Close()
	}()
}
