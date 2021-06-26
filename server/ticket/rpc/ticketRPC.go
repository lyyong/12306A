package rpc

import (
	"common/tools/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	pb "rpc/ticket/proto/ticketRPC"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models"
	"ticket/service"
	"ticket/utils/redispool"
	"time"
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
				Id:                uint32(res[j].ID),
				TrainNum:          res[j].TrainNum,
				StartStation:      res[j].StartStation,
				StartTime:         res[j].StartTime.String(),
				DestStation:       res[j].DestStation,
				DestTime:          res[j].DestTime.String(),
				SeatType:          res[j].SeatType,
				CarriageNumber:    res[j].CarriageNumber,
				SeatNumber:        res[j].SeatNumber,
				Price:             res[j].Price,
				CertificateNumber: res[j].CertificateNumber,
				PassengerName:     res[j].PassengerName,
				OrderOutsideId:    res[j].OrderOutsideId,
			}
		}
		list[i] = &pb.Tickets{Tickets: tickets}
	}

	return &pb.TicketsList{
		List: list,
	}, nil
}

func (ts *TicketServer) GetUnHandleTickets(ctx context.Context, in *pb.GetUnHandleTicketsRequest) (*pb.Tickets, error) {
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
			TrainId:           tpTickets[i].TrainId,
			TrainNum:          tpTickets[i].TrainNum,
			StartStationId:    tpTickets[i].StartStationId,
			StartStation:      tpTickets[i].StartStation,
			StartTime:         tpTickets[i].StartTime,
			DestStationId:     tpTickets[i].DestStationId,
			DestStation:       tpTickets[i].DestStation,
			DestTime:          tpTickets[i].ArriveTime,
			SeatType:          tpTickets[i].SeatType,
			CarriageNumber:    tpTickets[i].CarriageNumber,
			SeatNumber:        tpTickets[i].SeatNumber,
			Price:             tpTickets[i].Price,
			PassengerId:       tpTickets[i].PassengerId,
			PassengerName:     tpTickets[i].PassengerName,
			OrderOutsideId:    tpTickets[i].OrderId,
			CertificateNumber: tpTickets[i].CertificateNumber,
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
			Id:                uint32(res[j].ID),
			TrainNum:          res[j].TrainNum,
			StartStation:      res[j].StartStation,
			StartTime:         res[j].StartTime.String(),
			DestStation:       res[j].DestStation,
			DestTime:          res[j].DestTime.String(),
			SeatType:          res[j].SeatType,
			CarriageNumber:    res[j].CarriageNumber,
			SeatNumber:        res[j].SeatNumber,
			Price:             res[j].Price,
			PassengerName:     res[j].PassengerName,
			OrderOutsideId:    res[j].OrderOutsideId,
			CertificateNumber: res[j].CertificateNumber,
		}
	}
	return &pb.Tickets{Tickets: tickets}, nil
}

func (ts *TicketServer) UpdateTicketsState(ctx context.Context, in *pb.UpdateStateRequest) (*pb.Empty, error) {
	// 根据 ticketId 修改 Ticket 中对应记录，将 state 改为 in.State
	err := models.UpdateState(in.TicketsId, uint8(in.State))
	return &pb.Empty{}, err
}

func (ts *TicketServer) BuyTickets(ctx context.Context, in *pb.BuyTicketsRequest) (*pb.BuyTicketsResponseList, error) {
	// 处理请求数据
	passengers := make([]*ticketPoolPb.PassengerInfo, len(in.Passengers))
	for index, value := range in.Passengers {
		passengers[index] = &ticketPoolPb.PassengerInfo{
			PassengerName:     value.PassengerName,
			PassengerId:       value.PassengerId,
			CertificateNumber: value.CertificateNumber,
			SeatTypeId:        value.SeatTypeId,
		}
	}
	getTicketReq := &ticketPoolPb.GetTicketRequest{
		TrainId:        in.TrainId,
		StartStationId: in.StartStationId,
		DestStationId:  in.DestStationId,
		Date:           in.Date,
		Passengers:     passengers,
	}
	// rpc请求票池出票
	tpTickets, err := service.GetTickets(getTicketReq)
	if err != nil || len(tpTickets) == 0 {
		logging.Error(err)
		return nil, err
	}
	// 处理票池返回数据
	tickets := make([]models.Ticket, len(tpTickets))
	for i := 0; i < len(tpTickets); i++ {
		startTime, _ := time.Parse("2006-01-02 15:04", tpTickets[i].StartTime)
		arriveTime, _ := time.Parse("2006-01-02 15:04", tpTickets[i].ArriveTime)

		tickets[i] = models.Ticket{
			Model:             gorm.Model{},
			UserId:            in.UserId,
			TrainId:           tpTickets[i].TrainId,
			TrainNum:          tpTickets[i].TrainNum,
			StartStationId:    tpTickets[i].StartStationId,
			StartStation:      tpTickets[i].StartStation,
			StartTime:         startTime,
			DestStationId:     tpTickets[i].DestStationId,
			DestStation:       tpTickets[i].DestStation,
			DestTime:          arriveTime,
			SeatType:          tpTickets[i].SeatType,
			CarriageNumber:    tpTickets[i].CarriageNumber,
			SeatNumber:        tpTickets[i].SeatNumber,
			Price:             tpTickets[i].Price,
			OrderOutsideId:    in.OrderOuterId,
			PassengerId:       tpTickets[i].PassengerId,
			PassengerName:     tpTickets[i].PassengerName,
			CertificateNumber: tpTickets[i].CertificateNumber,
			State:             4,
		}
	}
	// 车票信息写入数据库
	err = models.AddMultipleTicket(&tickets)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	// 构造响应数据
	response := make([]*pb.BuyTicketsResponse, len(tickets))
	for i := 0; i < len(response); i++ {
		response[i] = &pb.BuyTicketsResponse{
			PassengerId: tickets[i].PassengerId,
			TicketId:    uint32(tickets[i].ID),
		}
	}
	return &pb.BuyTicketsResponseList{Response: response}, nil
}
