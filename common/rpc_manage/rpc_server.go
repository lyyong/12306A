// @Author LiuYong
// @Created at 2021-01-31
// @Modified at 2021-01-31
package rpc_manage

import (
	"common/router_tracer"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
)

// NewGRPCServer 创建一个Rpc服务器.
// opt... 服务器可选项, 例如过滤器.
// 如果已经开启了链路追踪怎会自定添加该功能
func NewGRPCServer(opt ...grpc.ServerOption) *grpc.Server {
	// 检查是否有开启链路追踪,
	// 如果有执行过 zipkin.Setup 方法则开启链路追踪
	if !router_tracer.IsTracing() {
		return grpc.NewServer(opt...)
	}
	zpClient, _ := router_tracer.GetClient()
	// return grpc.NewServer(grpc_middleware.WithUnaryServerChain(
	// 	// 链路追踪拦截器
	// 	otgrpc.OpenTracingServerInterceptor(*zpClient.Tracer(), otgrpc.LogPayloads()),
	// ))
	return grpc.NewServer(grpc.StatsHandler(zipkingrpc.NewServerHandler(zpClient.Tracer())))
}
