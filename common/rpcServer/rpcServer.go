// @Author LiuYong
// @Created at 2021-01-31
// @Modified at 2021-01-31
package rpcServer

import (
	"common/rpcServer/consul"
	"common/rpcServer/zipkin"
	"common/tools/logging"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
	"strconv"
)

type RpcServer struct {
	consulClient *consul.Client
	zipkinClient *zipkin.Client
	server       *grpc.Server
}

func (r RpcServer) Server() *grpc.Server {
	return r.server
}

// NewRpcServer 创建一个Rpc服务器
// opt... 服务器可选项, 例如过滤器
func NewRpcServer(opt ...grpc.ServerOption) *RpcServer {
	return &RpcServer{
		consulClient: nil,
		zipkinClient: nil,
		server:       grpc.NewServer(opt...),
	}
}

// NewRpcServerWithServerFind 创建一个带有服务发现的Rpc服务器
// name - 服务名称
// host - 服务的地址
// port - 服务的端口号
// serviceID - 服务的ID, 一般是name-host-port的形式
// target - consul的地址,例如"localhost:8500"
// interval - //consul的地址,例如"localhost:8500"
// ttl - 注册信息的缓存时间, 如果ttl过时前没有得到updateTTL则注册服务信息将被抛弃, 单位秒
// opt... 服务器可选项, 例如过滤器
func NewRpcServerWithServerFind(name string, host string, port int, serviceID string, target string, interval int, ttl int, opt ...grpc.ServerOption) (*RpcServer, error) {
	consulClient, err := consul.NewClient(name, host, port, serviceID, target, interval, ttl)
	if err != nil {
		return nil, err
	}
	if err := consulClient.Register(); err != nil {
		return nil, err
	}
	logging.Info("启动服务发现, serviceID:", serviceID, " consul服务地址:", target)
	return &RpcServer{
		consulClient: consulClient,
		zipkinClient: nil,
		server:       grpc.NewServer(opt...),
	}, nil
}

// NewRpcServerWithServerFindAndHttpTracer 创建一个带有服务发现的Rpc服务器
// name - 服务名称
// host - 服务的地址
// port - 服务的端口号
// serviceID - 服务的ID, 一般是name-host-port的形式
// target - consul的地址,例如"localhost:8500"
// interval - //consul的地址,例如"localhost:8500"
// ttl - 注册信息的缓存时间, 如果ttl过时前没有得到updateTTL则注册服务信息将被抛弃, 单位秒
// zipkinHttp - zipkin所在的http接受地址 例如"http://localhost:9411/api/v2/spans"
// opt... 服务器可选项, 例如过滤器
func NewRpcServerWithServerFindAndHttpTracer(name string, host string, port int, serviceID string, target string, interval int, ttl int, zipkinHttp string, opt ...grpc.ServerOption) (*RpcServer, error) {
	consulClient, err := consul.NewClient(name, host, port, serviceID, target, interval, ttl)
	if err != nil {
		return nil, err
	}
	if err := consulClient.Register(); err != nil {
		return nil, err
	}
	logging.Info("启动服务发现, serviceID:", serviceID, " consul服务地址:", target)
	zipkinClient := zipkin.NewClient(name, host, strconv.Itoa(port))
	err = zipkinClient.SetupByHttp(zipkinHttp)
	if err != nil {
		return nil, err
	}
	logging.Info("启动链路追踪, serviceID:", serviceID, " zipkin服务地址:", zipkinHttp)
	return &RpcServer{
		consulClient: consulClient,
		zipkinClient: nil,
		// 创建grpc服务并且设置一元中间件
		server: grpc.NewServer(grpc_middleware.WithUnaryServerChain(
			// 链路追踪拦截器
			otgrpc.OpenTracingServerInterceptor(*zipkinClient.Tracer(), otgrpc.LogPayloads()),
		)),
	}, nil
}

func (r *RpcServer) Close() {
	if r.zipkinClient != nil {
		r.zipkinClient.Close()
	}
	if r.consulClient != nil {
		r.consulClient.DeRegister()
	}
}
