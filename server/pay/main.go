package main

import (
	"12306A/server/pay/router"
	"12306A/server/pay/tools/logging"
	"12306A/server/pay/tools/setting"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// 需要初始化的组件
func init() {
	setting.Setup()

	logging.Setup()
}

// 需要关闭的组件
func serverClose() {

}

func main() {
	r := router.InitRouter()

	logging.Info("启动Pay服务, 端口号: ", setting.ServerSetting.HttpPort)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:      r,
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
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logging.Fatal("Server Shutdown: ", err)
	}

	logging.Info("Server exiting")
	serverClose()
}