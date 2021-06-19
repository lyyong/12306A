// @Author: KongLingWen
// @Created at 2021/2/22
// @Modified at 2021/2/22

package client

import (
	"common/rpc_manage"
	"context"
	"errors"
	"google.golang.org/grpc"
	"rpc/ticket/proto/ticketRPC"
	"sync"
)

type TicketRPCClient struct {
	ticketClient *ticketRPC.TicketServiceClient
}

var (
	client *TicketRPCClient
	once sync.Once
)

const targetServiceName = ":9442"

// NewClient 创建一个ticketPool的RPC客户端
func NewClient() (*TicketRPCClient, error) {
	return NewClientWithTarget(targetServiceName)
}

func NewClientWithTarget(target string) (*TicketRPCClient, error) {
	var err error
	once.Do(func() {
		client = &TicketRPCClient{ticketClient: nil}
		var conn *grpc.ClientConn
		conn, err = rpc_manage.NewGRPCClientConn(target)
		if err != nil {
			client = nil
			return
		}
		tclient := ticketRPC.NewTicketServiceClient(conn)
		client.ticketClient = &tclient
	})

	return client, err
}

func(c TicketRPCClient) GetTicketByOrdersId(request *ticketRPC.GetTicketByOrdersIdRequest)(*ticketRPC.TicketsList, error){
	if c.ticketClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.ticketClient
	return tclient.GetTicketByOrdersId(context.Background(), request)
}

func(c TicketRPCClient) GetTicketByPassengerId(request *ticketRPC.GetTicketByPassengerIdRequest)(*ticketRPC.Tickets, error){
	if c.ticketClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.ticketClient
	return tclient.GetTicketByPassengerId(context.Background(), request)
}

func(c TicketRPCClient) GetUnHandleTickets(request *ticketRPC.GetUnHandleTicketsRequest)(*ticketRPC.Tickets, error){
	if c.ticketClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.ticketClient
	return tclient.GetUnHandleTickets(context.Background(), request)
}

func(c TicketRPCClient) BuyTickets(request *ticketRPC.BuyTicketsRequest)(*ticketRPC.BuyTicketsResponseList, error){
	if c.ticketClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.ticketClient
	return tclient.BuyTickets(context.Background(), request)
}

func(c TicketRPCClient) UpdateTicketsState(request *ticketRPC.UpdateStateRequest)(*ticketRPC.Empty, error) {
	if c.ticketClient == nil {
		return nil, errors.New("没有NewClient")
	}
	tclient := *c.ticketClient
	return tclient.UpdateTicketsState(context.Background(), request)
}