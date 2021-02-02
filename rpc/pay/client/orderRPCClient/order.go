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

const targetServiceName = "pay-server"

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
