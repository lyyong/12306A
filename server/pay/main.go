// @Author liuYong
// @Created at 2020-01-05
package main

import (
	"common/router_tracer"
	"common/server_find"
	"common/tools/logging"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"os"
	"os/signal"
	"pay/router"
	"pay/tools/cache"
	"pay/tools/database"
	"pay/tools/setting"
	"strconv"
	"strings"
	"time"
)

// 需要初始化的组件
func init() {
	// 载入配置文件
	setting.Setup()
	// 加载日子系统
	if setting.Server.RunMode == "debug" {
		logging.SetupWithMode(logging.LogDebug)
	} else {
		logging.SetupWithMode(logging.LogRelease)
	}
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

	// 加载redis
	err = cache.Setup(&redis.Options{
		Network: "tcp",
		Addr:    setting.Redis.Host,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			_, err := cn.Ping(ctx).Result()
			if err != nil {
				logging.Error(err)
				return err
			}
			// logging.Info(res)
			return nil
		},
		ReadTimeout:  setting.Redis.ReadTimeout,
		WriteTimeout: setting.Redis.WriteTimeout,
		PoolSize:     setting.Redis.PoolSize,
		MinIdleConns: setting.Redis.MinIdleConns,
		IdleTimeout:  setting.Redis.IdelTimeout,
	})
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
	// gin 路由
	ginRouter := router.InitRouter()
	// grpc路由
	grpcServer := router.InitRPCService()

	logging.Info("启动Pay服务, 端口号: ", setting.Server.HttpPort)
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", setting.Server.HttpPort),
		// 使用h2c库,避免使用TLS, 实现http2的未加密模式
		Handler: h2c.NewHandler(http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				if request.ProtoMajor == 2 && strings.Contains(
					request.Header.Get("Content-Type"), "application/grpc") {
					// grpc请求
					grpcServer.ServeHTTP(responseWriter, request)
				} else {
					ginRouter.ServeHTTP(responseWriter, request)
				}
				return
			}), &http2.Server{}),
		ReadTimeout:  setting.Server.ReadTimeout,
		WriteTimeout: setting.Server.WriteTimeout,
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
