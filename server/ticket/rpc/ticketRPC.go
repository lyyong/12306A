package rpc

import (
	"context"
	pb "rpc/ticket/proto/ticketRPC"
)

type TicketServer struct {
	/*
		rpc2 AddTickets (Tickets) returns (Empty){}
	  	rpc2 GetTicketByIndentId (GetByIndentRequest) returns (Tickets){}
	  	rpc2 GetTicketByPassengerId (GetByPassengerRequest) returns (Tickets){}
	  	rpc2 UpdateState (UpdateStateRequest) returns (Empty){}
	*/
}



func (ts *TicketServer) AddTickets(ctx context.Context, in *pb.Tickets) (*pb.Empty, error) {
	// 生成外部订单号

	// 写入redis，设置过期时间

	// 返回外部订单号

	return &pb.Empty{}, nil
}

func (ts *TicketServer) GetTicketByIndentId(ctx context.Context, in *pb.GetByIndentRequest) (*pb.Tickets, error) {
	// 根据订单号查询 Ticket 表

	return &pb.Tickets{}, nil
}

func (ts *TicketServer) GetTicketByPassengerId(ctx context.Context, in *pb.GetByPassengerRequest) (*pb.Tickets, error) {
	// 根据 passengerId 查询 Ticket 表

	return &pb.Tickets{}, nil
}

func (ts *TicketServer) UpdateState(ctx context.Context, in *pb.UpdateStateRequest) (*pb.Empty, error) {
	// 根据 ticketId 修改 Ticket 中对应记录，将 state 改为 in.State

	return &pb.Empty{}, nil
}