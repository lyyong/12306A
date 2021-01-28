// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package main

import (
	"context"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"os"
	"os/signal"
	"pay/router"
	"pay/script"
	"pay/tools/logging"
	"pay/tools/setting"
	"strings"
	"time"
)

// 需要初始化的组件
func init() {
	logging.Setup()
	setting.Setup()
	script.Setup()
}

// 需要关闭的组件
func serverClose() {

}

func main() {
	// gin 路由
	ginRouter := router.InitRouter()
	// grpc路由
	grpcRouter := router.InitRPCService()
	logging.Info("启动Pay服务, 端口号: ", setting.ServerSetting.HttpPort)
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		// 使用h2c库,避免使用TLS, 实现http2的未加密模式
		Handler: h2c.NewHandler(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				if request.ProtoMajor == 2 && strings.Contains(
					request.Header.Get("Content-Type"), "application/grpc") {
					// grpc请求
					grpcRouter.ServeHTTP(responseWriter, request)
				} else {
					ginRouter.ServeHTTP(responseWriter, request)
				}
				return
			}), &http2.Server{}),
		ReadTimeout:  setting.ServerSetting.ReadTimeout,
		WriteTimeout: setting.ServerSetting.WriteTimeout,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logging.Error("Listens: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.Info("Shutdown Server ......")

	// 5 秒后关闭资源
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		serverClose()
	}()
	if err := s.Shutdown(ctx); err != nil {
		logging.Fatal("Server Shutdown: ", err)
	}

	logging.Info("Server exiting")
}
