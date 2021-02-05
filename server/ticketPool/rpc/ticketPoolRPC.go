package rpc

import (
	"context"
	pb "rpc/ticketPool/proto/ticketPoolRPC"
	"strconv"
)

type TicketPoolServer struct {
	/*
		rpc2 GetTicket (GetTicketRequest) returns (GetTicketResponse){}
	  	rpc2 GetTicketNumber (GetTicketNumberRequest) returns (GetTicketNumberResponse){}
	  	rpc2 RefundTicket (RefundTicketRequest) returns (RefundTicketResponse){}
	*/
}

func (tps *TicketPoolServer) GetTicket(ctx context.Context, in *pb.GetTicketRequest) (*pb.GetTicketResponse, error) {
	// 根据请求数据出票

	ticketNumber := len(in.Passengers)
	tickets := make([]*pb.Ticket, ticketNumber)
	for index, value := range in.Passengers {
		tickets[index] = &pb.Ticket{
			Id:             int32(index),
			TrainId:        in.TrainId,
			StartStationId: in.StartStationId,
			StartTime:      "10:37",
			DestStationId:  in.DestStationId,
			ArriveTime:     "14:56",
			Date:           in.Date,
			SeatTypeId:     value.SeatTypeId,
			CarriageNumber: "5",
			SeatNumber:     strconv.Itoa(index+20),
			PassengerId:    value.PassengerId,
			IndentId:       0,
			Amount:         2650,
		}

	}

	return &pb.GetTicketResponse{
		Tickets: tickets,
	}, nil
}

func (tps *TicketPoolServer) GetTicketNumber(ctx context.Context, in *pb.GetTicketNumberRequest) (*pb.GetTicketNumberResponse, error) {
	// 查询车次余票

	trainsTicketInfo := make([]*pb.TrainTicketInfo, 10)
	return &pb.GetTicketNumberResponse{
		TrainsTicketInfo: trainsTicketInfo,
	}, nil
}

func (tps *TicketPoolServer) RefundTicket(ctx context.Context, in *pb.RefundTicketRequest) (*pb.RefundTicketResponse, error) {
	// 退票到票池

	return &pb.RefundTicketResponse{
		IsOk: true,
	}, nil
}