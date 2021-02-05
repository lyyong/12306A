package indentRPCClient

import (
	"common/rpc_manage"
	"context"
	"rpc/indent/proto/indentRPC"
)



type IndentRPCClient struct {
	pbClient *indentRPC.IndentServiceClient
}

// 唯一
var client *IndentRPCClient

const targetServiceName = "indent-server"

// NewClient 创建一个 OrderRPCClient
func NewClient() (*IndentRPCClient, error) {
	if client != nil {
		return client, nil
	}
	client = &IndentRPCClient{nil}
	conn, err := rpc_manage.NewGRPCClientConn(targetServiceName)
	if err != nil {
		client = nil
		return nil, err
	}
	tclient := indentRPC.NewIndentServiceClient(conn)
	client.pbClient = &tclient
	return client, nil
}

func (c *IndentRPCClient) Create(req *indentRPC.CreateRequest) (*indentRPC.CreateResponse, error) {
	tclient := *c.pbClient
	resp, err := tclient.Create(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
