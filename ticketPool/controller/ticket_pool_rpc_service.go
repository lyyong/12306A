// @Author LiuYong
// @Created at 2021-02-15
// @Modified at 2021-02-15
package controller

import (
	"context"
	"rpc/ticketPool/proto/ticketPoolRPC"
)

type TicketPoolRPCService struct {
}

func (t TicketPoolRPCService) GetTicket(ctx context.Context, request *ticketPoolRPC.GetTicketRequest) (*ticketPoolRPC.GetTicketResponse, error) {
	panic("implement me")
}

func (t TicketPoolRPCService) GetTicketNumber(ctx context.Context, request *ticketPoolRPC.GetTicketNumberRequest) (*ticketPoolRPC.GetTicketNumberResponse, error) {
	panic("implement me")
}

func (t TicketPoolRPCService) RefundTicket(ctx context.Context, request *ticketPoolRPC.RefundTicketRequest) (*ticketPoolRPC.RefundTicketResponse, error) {
	panic("implement me")
}
