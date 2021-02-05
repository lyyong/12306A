// @Author LiuYong
// @Created at 2021-02-02
// @Modified at 2021-02-02
package main

import (
	"candidate/router"
	"candidate/tools/setting"
	"common/router_tracer"
	"common/server_find"
	"common/tools/logging"
	"context"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"os"
	"os/signal"
	"pay/tools/database"
	"strconv"
	"strings"
	"time"
)

// 需要初始化的组件
func init() {
	// 加载日子系统
	logging.Setup()
	// 载入配置文件
	setting.Setup()
	// 服务发现
	server_find.Register(setting.Server.Name,
		setting.Server.Host, strconv.Itoa(setting.Server.HttpPort), setting.Consul.ServiceID, setting.Consul.Address, setting.Consul.Interval, setting.Consul.TTL)
	// 链路追踪
	err := router_tracer.SetupByHttp(setting.Server.Name,
		setting.Server.Host, strconv.Itoa(setting.Server.HttpPort), setting.Zipkin.HttpEndpoint)
	if err != nil {
		logging.Error(err)
	}
	// 加载数据库
	err = database.Setup(setting.Database.Type, setting.Database.Username, setting.Database.Password, setting.Database.Host, setting.Database.DbName)
	if err != nil {
		logging.Error(err)
	}
}

// 需要关闭的组件
func serverClose() {
	server_find.DeRegister()
	router_tracer.Close()
	database.Close()
}

func main() {
	ginRouter := router.InitRouter()
	logging.Info("启动Candidate服务, 端口号: ", setting.Server.HttpPort)
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", setting.Server.HttpPort),
		Handler: h2c.NewHandler(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				if request.ProtoMajor == 2 &&
					strings.Contains(request.Header.Get("Content-type"), "application/grpc") {
					logging.Info("grpc 请求")
				} else {
					ginRouter.ServeHTTP(responseWriter, request)
				}
				return
			},
		), &http2.Server{}),
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logging.Error("Listens: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	logging.Info("Shutdown Server ......")

	// 5 秒后关闭资源
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		serverClose()
	}()
	if err := s.Shutdown(ctx); err != nil {
		logging.Error("Server Shutdown: ", err)
	}

	logging.Info("Server exiting")
}
