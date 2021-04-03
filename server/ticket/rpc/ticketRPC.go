package rpc

import (
	"context"
	pb "rpc/ticket/proto/ticketRPC"
	"ticket/models"
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

func (ts *TicketServer) GetTicketByOrdersId(ctx context.Context, in *pb.GetTicketByOrdersIdRequest) (*pb.TicketsList, error) {
	// 根据订单号查询 Ticket 表
	list := make([]*pb.Tickets, len(in.OrdersId))
	for i := 0; i < len(in.OrdersId); i++ {
		orderId := in.OrdersId[i]

		res, err := models.GetTicketByOrderId(orderId)
		if err != nil {
			return nil, err
		}
		tickets := make([]*pb.Ticket, len(res))
		for j := 0; j < len(res); j++ {
			tickets[j] = &pb.Ticket{
				TrainNum:       res[j].TrainNum,
				StartStation:   res[j].StartStation,
				StartTime:      res[j].StartTime.String(),
				DestStation:    res[j].DestStation,
				DestTime:       res[j].DestTime.String(),
				SeatType:       res[j].SeatType,
				CarriageNumber: res[j].CarriageNumber,
				SeatNumber:     res[j].SeatNumber,
				Price:          res[j].Price,
				PassengerName:  res[j].PassengerName,
				OrderOutsideId: res[j].OrderOutsideId,
			}
		}
		list[i] = &pb.Tickets{Tickets: tickets}
	}

	return &pb.TicketsList{
		List: list,
	}, nil
}

func (ts *TicketServer) GetTicketByPassengerId(ctx context.Context, in *pb.GetTicketByPassengerIdRequest) (*pb.Tickets, error) {
	// 根据 passengerId 查询 Ticket 表

	return &pb.Tickets{}, nil
}

func (ts *TicketServer) UpdateState(ctx context.Context, in *pb.UpdateStateRequest) (*pb.Empty, error) {
	// 根据 ticketId 修改 Ticket 中对应记录，将 state 改为 in.State

	return &pb.Empty{}, nil
}