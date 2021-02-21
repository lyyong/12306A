// @Author LiuYong
// @Created at 2021-02-17
// @Modified at 2021-02-17
package Client

import (
	"common/rpc_manage"
	"context"
	"errors"
	"rpc/ticketPool/proto/ticketPoolRPC"
)

type TPRPCClient struct {
	tpClient *ticketPoolRPC.TicketPoolServiceClient
}

var client *TPRPCClient

const targetServiceName = "ticketPool"

// NewClient 创建一个ticketPool的RPC客户端
func NewClient() (*TPRPCClient, error) {
	if client != nil {
		return client, nil
	}
	client = &TPRPCClient{tpClient: nil}
	conn, err := rpc_manage.NewGRPCClientConn(targetServiceName)
	if err != nil {
		client = nil
		return nil, err
	}
	tclient := ticketPoolRPC.NewTicketPoolServiceClient(conn)
	client.tpClient = &tclient
	return client, nil
}

func (c TPRPCClient) GetTicket(request *ticketPoolRPC.GetTicketRequest) (*ticketPoolRPC.GetTicketResponse, error) {
	if c.tpClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.tpClient
	return tclient.GetTicket(context.Background(), request)
}

func (c TPRPCClient) GetTicketNumber(request *ticketPoolRPC.GetTicketNumberRequest) (*ticketPoolRPC.GetTicketNumberResponse, error) {
	if c.tpClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.tpClient
	return tclient.GetTicketNumber(context.Background(), request)
}

func (c TPRPCClient) RefundTicket(request *ticketPoolRPC.RefundTicketRequest) (*ticketPoolRPC.RefundTicketResponse, error) {
	if c.tpClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.tpClient
	return tclient.RefundTicket(context.Background(), request)
}
