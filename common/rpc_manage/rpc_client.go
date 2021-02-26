// @Author LiuYong
// @Created at 2021-01-31
// @Modified at 2021-01-31
package rpc_manage

import (
	"common/router_tracer"
	"common/server_find"
	"fmt"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"strings"
)

const (
	// https://github.com/grpc/grpc/blob/master/doc/naming.md 服务配置详细地址
	SERVICE_CONFIG = `{
	  "loadBalancingConfig": [ { "%s": {} } ]
	}`

	ROUND_ROBIN = "round_robin" // 轮询负载均衡
	GRPCLB      = "grpclb"
)

// NewGRPCClientConn 创建grpc客户端连接,
// 如果开启了服务注册功能则会自动开启负载均衡, 暂时只提供轮询
// 如果开启了链路追踪功能则会自动加上,
// targetService - 目标服务如果没有开启服务发现, 可以直接是host, 如果开启了服务发现最好直接使用目标的服务名, 使用包含了":"号就是
// 使用端口的直连模式, 或者nginx负载模式
// targetService 值为nginx:port将启动nginx负载均衡,  这里的port是nginx要监听的端口号
func NewGRPCClientConn(targetService string) (*grpc.ClientConn, error) {

	if router_tracer.IsTracing() && server_find.IsRegister() && !strings.Contains(targetService, ":") {
		// 开启链路追踪和服务注册, 目标服务不包含:表示使用的服务名
		zkClient, _ := router_tracer.GetClient()
		cClient, _ := server_find.GetClient()
		disc, err := server_find.NewServiceDiscoveryAboutBalance(cClient, targetService)
		if err != nil {
			return nil, err
		}
		// 注册Builder
		resolver.Register(disc)
		return grpc.Dial(server_find.SCHEME+":///",
			grpc.WithStatsHandler(zipkingrpc.NewClientHandler(zkClient.Tracer())),
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(SERVICE_CONFIG, ROUND_ROBIN)))
	}

	if router_tracer.IsTracing() {
		// 开启了链路追踪
		zkClient, _ := router_tracer.GetClient()
		// return grpc.Dial(targetService,
		// 	grpc.WithUnaryInterceptor(
		// 		otgrpc.OpenTracingClientInterceptor(*zkClient.Tracer(), otgrpc.LogPayloads())),
		// 	grpc.WithInsecure())
		return grpc.Dial(targetService,
			grpc.WithStatsHandler(zipkingrpc.NewClientHandler(zkClient.Tracer())),
			grpc.WithInsecure())
	}

	if server_find.IsRegister() && !strings.Contains(targetService, ":") {
		// 开启了服务发现
		cClient, _ := server_find.GetClient()
		disc, err := server_find.NewServiceDiscoveryAboutBalance(cClient, targetService)
		if err != nil {
			return nil, err
		}
		// 注册Builder
		resolver.Register(disc)
		return grpc.Dial(server_find.SCHEME+":///",
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(SERVICE_CONFIG, ROUND_ROBIN)))
	}
	// 都没有开启
	return grpc.Dial(targetService, grpc.WithInsecure())
}
