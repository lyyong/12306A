// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCClient

import (
	"common/rpc_manage"
	"context"
	"rpc/pay/proto/orderRPCpb"
)

type OrderRPCClient struct {
	pbClient *orderRPCpb.OrderRPCServiceClient
}

// 唯一
var client *OrderRPCClient

const targetServiceName = "nginx:18082"

// NewClient 创建一个 OrderRPCClient
func NewClient() (*OrderRPCClient, error) {
	if client != nil {
		return client, nil
	}
	client = &OrderRPCClient{nil}
	conn, err := rpc_manage.NewGRPCClientConn(targetServiceName)
	if err != nil {
		client = nil
		return nil, err
	}
	tclient := orderRPCpb.NewOrderRPCServiceClient(conn)
	client.pbClient = &tclient
	return client, nil
}

func (c *OrderRPCClient) Create(info *orderRPCpb.CreateRequest) (*orderRPCpb.CreateRespond, error) {
	tclient := *c.pbClient
	resp, err := tclient.Create(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *OrderRPCClient) Read(info *orderRPCpb.SearchCondition) (*orderRPCpb.ReadRespond, error) {
	tclient := *c.pbClient
	resp, err := tclient.Read(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *OrderRPCClient) UpdateState(info *orderRPCpb.UpdateStateRequest) (*orderRPCpb.Respond, error) {
	tclient := *c.pbClient
	resp, err := tclient.UpdateState(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *OrderRPCClient) UpdateStateWithRelativeOrder(info *orderRPCpb.UpdateStateWithRRequest) (*orderRPCpb.Respond, error) {
	tclient := *c.pbClient
	resp, err := tclient.UpdateStateWithRelativeOrder(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *OrderRPCClient) GetNoFinishOrder(info *orderRPCpb.SearchCondition) (*orderRPCpb.OrderInfo, error) {
	tclient := *c.pbClient
	resp, err := tclient.GetNoFinishOrder(context.Background(), info)
	if err != nil {
		return nil, err
	}
	return resp, err
}
