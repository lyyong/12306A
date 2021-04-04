package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	pb "rpc/ticket/proto/ticketRPC"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models"
	"ticket/utils/redispool"
)

type TicketServer struct {
	/*
		rpc2 AddTickets (Tickets) returns (Empty){}
	  	rpc2 GetTicketByIndentId (GetByIndentRequest) returns (Tickets){}
	  	rpc2 GetTicketByPassengerId (GetByPassengerRequest) returns (Tickets){}
	  	rpc2 UpdateState (UpdateStateRequest) returns (Empty){}
	*/
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

func (ts *TicketServer) GetUnHandleTickets(ctx context.Context, in *pb.GetUnHandleTicketsRequest) (*pb.Tickets, error){
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("ticket_%d", in.UserId)

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	var tpTickets []*ticketPoolPb.Ticket
	err = json.Unmarshal(data, &tpTickets)
	if err != nil {
		return nil, err
	}

	tickets := make([]*pb.Ticket, len(tpTickets))
	for i := 0; i < len(tpTickets); i++ {

		tickets[i] = &pb.Ticket{
			TrainNum:       tpTickets[i].TrainNum,
			StartStation:   tpTickets[i].StartStation,
			StartTime:      tpTickets[i].StartTime,
			DestStation:    tpTickets[i].DestStation,
			DestTime:       tpTickets[i].ArriveTime,
			SeatType:       tpTickets[i].SeatType,
			CarriageNumber: tpTickets[i].CarriageNumber,
			SeatNumber:     tpTickets[i].SeatNumber,
			Price:          tpTickets[i].Price,
			PassengerName:  tpTickets[i].PassengerName,
			OrderOutsideId: "",
		}
	}
	return &pb.Tickets{Tickets: tickets}, nil
}

func (ts *TicketServer) GetTicketByPassengerId(ctx context.Context, in *pb.GetTicketByPassengerIdRequest) (*pb.Tickets, error) {
	// 根据 passengerId 查询 Ticket 表
	res, err := models.GetTicketsByPassengerId(in.PassengerId)
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
	return &pb.Tickets{Tickets: tickets}, nil
}

func (ts *TicketServer) UpdateState(ctx context.Context, in *pb.UpdateStateRequest) (*pb.Empty, error) {
	// 根据 ticketId 修改 Ticket 中对应记录，将 state 改为 in.State

	return &pb.Empty{}, nil
}