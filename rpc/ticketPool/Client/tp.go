// @Author LiuYong
// @Created at 2021-02-17
// @Modified at 2021-02-17
package Client

import (
	"common/rpc_manage"
	"context"
	"errors"
	"google.golang.org/grpc"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"sync"
)

type TPRPCClient struct {
	tpClient *ticketPoolRPC.TicketPoolServiceClient
}

var (
	client *TPRPCClient
	once   sync.Once
)

const targetServiceName = ":9443"

// NewClient 创建一个ticketPool的RPC客户端
func NewClient() (*TPRPCClient, error) {
	return NewClientWithTarget(targetServiceName)
}

// NewClientWithTarget 创建一个ticketPool的RPC客户端
func NewClientWithTarget(target string) (*TPRPCClient, error) {
	var err error
	once.Do(func() {
		client = &TPRPCClient{tpClient: nil}
		var conn *grpc.ClientConn
		conn, err = rpc_manage.NewGRPCClientConn(target)
		if err != nil {
			client = nil
			return
		}
		tclient := ticketPoolRPC.NewTicketPoolServiceClient(conn)
		client.tpClient = &tclient
	})

	return client, err
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
