// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCClient

import (
	"common/rpcServer/zipkin"
	"context"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
	"rpc/pay/proto/orderRPCpb"
)

type OrderRPCClient struct {
	pbClient *orderRPCpb.OrderRPCServiceClient
	conn     *grpc.ClientConn
	zpClient *zipkin.Client
}

var client *OrderRPCClient

const server_port = "8082"

// NewClient 创建一个不带链路追踪的客户端
func NewClient() (*OrderRPCClient, error) {
	if client != nil {
		return client, nil
	}
	client = &OrderRPCClient{
		pbClient: nil,
		conn:     nil,
		zpClient: nil,
	}
	var err error

	// Dial的功能是连接一个服务器, 服务器可以提供grpc服务端的信息
	client.conn, err = grpc.Dial(":"+server_port, grpc.WithInsecure())
	if err != nil {
		client = nil
		return nil, err
	}
	tclient := orderRPCpb.NewOrderRPCServiceClient(client.conn)
	client.pbClient = &tclient
	return client, nil
}

// NewClientWithHttpTracer 创建一个带有http链路追踪的客户端
// servicenName - 服务名称
// serviceHost - 服务的地址
// servicePort - 服务的端口号
// zipkinHttp - zipkin所在的http接受地址 例如"http://localhost:9411/api/v2/spans"
func NewClientWithHttpTracer(servicenName string, serviceHost string, servicePort string, zipkinHttp string) (*OrderRPCClient, error) {
	if client != nil && client.zpClient != nil {
		return client, nil
	}
	client = &OrderRPCClient{
		pbClient: nil,
		conn:     nil,
		zpClient: nil,
	}
	var err error
	cli := zipkin.NewClient(servicenName, serviceHost, servicePort)
	err = cli.SetupByHttp(zipkinHttp)
	if err != nil { // 创建zipkin追踪器出错
		client = nil
		return nil, err
	}
	// 设置zp客户端
	client.zpClient = cli
	// 设置zp客户端
	client.zpClient = cli
	client.conn, err = grpc.Dial(":"+server_port,
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(*cli.Tracer(), otgrpc.LogPayloads())),
		grpc.WithInsecure())
	if err != nil {
		client = nil
		return nil, err
	}
	tclient := orderRPCpb.NewOrderRPCServiceClient(client.conn)
	client.pbClient = &tclient
	return client, nil
}

func (c *OrderRPCClient) Close() error {
	//zipkin.Close()
	if c.conn == nil {
		return nil
	}
	c.zpClient.Close()
	defer func() {
		c.conn = nil
		client = nil
	}()
	return c.conn.Close()
}

func (c *OrderRPCClient) Create(info *orderRPCpb.CreateInfo) (*orderRPCpb.Error, error) {
	tclient := *c.pbClient
	resp, err := tclient.Create(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *OrderRPCClient) Read(info *orderRPCpb.SearchInfo) (*orderRPCpb.Info, error) {
	tclient := *c.pbClient
	resp, err := tclient.Read(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *OrderRPCClient) UpdateState(info *orderRPCpb.UpdateStateInfo) (*orderRPCpb.Error, error) {
	tclient := *c.pbClient
	resp, err := tclient.UpdateState(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *OrderRPCClient) UpdateStateWithRelativeOrder(info *orderRPCpb.UpdateStateWithRInfo) (*orderRPCpb.Error, error) {
	tclient := *c.pbClient
	resp, err := tclient.UpdateStateWithRelativeOrder(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}
