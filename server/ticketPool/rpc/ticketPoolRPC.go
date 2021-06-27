package rpc

import (
	"common/tools/logging"
	"context"
	"fmt"
	pb "rpc/ticketPool/proto/ticketPoolRPC"
	"strings"
	"ticketPool/ticketpool"
	"time"
)

type TicketPoolServer struct {
	/*
			rpc GetTicket (GetTicketRequest) returns (GetTicketResponse){}
		  	rpc GetTicketNumber (GetTicketNumberRequest) returns (GetTicketNumberResponse){}
		  	rpc RefundTicket (RefundTicketRequest) returns (RefundTicketResponse){}
	*/
}

func (tps *TicketPoolServer) GetTicket(ctx context.Context, req *pb.GetTicketRequest) (*pb.GetTicketResponse, error) {
	logging.Info("购票请求: ", req)

	tp := ticketpool.Tp

	seatCountMap := make(map[uint32]int32)
	for i := 0; i < len(req.Passengers); i++ {
		seatTypeId := req.Passengers[i].SeatTypeId
		seatCountMap[seatTypeId]++
	}
	// 调用tp.GetTicket出票，得到的是座位切片
	seatsMap, err := tp.GetTicket(req.TrainId, req.StartStationId, req.DestStationId, req.Date, seatCountMap)
	if err != nil {
		logging.Error(err)
		return &pb.GetTicketResponse{Tickets: nil}, err
	}

	// 整合票池返回的座位信息与请求信息，生成车票返回
	startStation := tp.GetStopStation(req.TrainId, req.StartStationId)
	destStation := tp.GetStopStation(req.TrainId, req.DestStationId)
	st, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", req.Date, startStation.StartTime), time.Local)
	at, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", req.Date, destStation.ArriveTime), time.Local)
	if at.Before(st) {
		at.AddDate(0, 0, 1)
	}
	startTime := st.Format("2006-01-02 15:04")
	arriveTime := at.Format("2006-01-02 15:04")

	tickets := make([]*pb.Ticket, len(req.Passengers))
	ticketIndex := 0
	for seatTypeId, seats := range seatsMap {
		seatIndex := 0
		seatType := tp.GetSeatTypeNameById(seatTypeId)
		for i := 0; i < len(req.Passengers); i++ {
			passengerInfo := req.Passengers[i]
			if seatTypeId == passengerInfo.SeatTypeId {
				carriageAndSeat := strings.Split(seats[seatIndex], " ")
				// 生成车票信息
				tickets[ticketIndex] = &pb.Ticket{
					Id:                0,
					TrainId:           req.TrainId,
					TrainNum:          tp.GetTrainNumber(req.TrainId),
					StartStationId:    req.StartStationId,
					StartStation:      startStation.StationName,
					StartTime:         startTime,
					DestStationId:     req.DestStationId,
					DestStation:       destStation.StationName,
					ArriveTime:        arriveTime,
					SeatTypeId:        seatTypeId,
					SeatType:          seatType,
					CarriageNumber:    carriageAndSeat[0],
					SeatNumber:        carriageAndSeat[1],
					PassengerName:     passengerInfo.PassengerName,
					PassengerId:       passengerInfo.PassengerId,
					CertificateNumber: passengerInfo.CertificateNumber,
					OrderId:           "",
					Price:             88,
				}
				ticketIndex++
				seatIndex++
			}
		}
	}
	return &pb.GetTicketResponse{Tickets: tickets}, nil
}

func (tps *TicketPoolServer) GetTicketNumber(ctx context.Context, req *pb.GetTicketNumberRequest) (*pb.GetTicketNumberResponse, error) {
	logging.Info("查询余票请求: ", req)

	// 查询车次余票
	tp := ticketpool.Tp
	condition := req.Condition
	tti := make([]*pb.TrainTicketInfo, len(condition))

	for i := 0; i < len(condition); i++ {
		seatCountMap, err := tp.SearchTicketCount(condition[i].TrainId, condition[i].StartStationId, condition[i].DestStationId, req.Date)
		if err != nil {
			logging.Error(err)
			return nil, err
		}
		seatInfo := make([]*pb.SeatInfo, len(seatCountMap))
		index := 0
		for seatTypeId, count := range seatCountMap {
			seatInfo[index] = &pb.SeatInfo{SeatTypeId: seatTypeId, SeatNumber: count}
			index++
		}
		tti[i] = &pb.TrainTicketInfo{TrainId: condition[i].TrainId, SeatInfo: seatInfo}
	}

	return &pb.GetTicketNumberResponse{TrainsTicketInfo: tti}, nil
}

func (tps *TicketPoolServer) RefundTicket(ctx context.Context, req *pb.RefundTicketRequest) (*pb.RefundTicketResponse, error) {
	// 退票到票池
	logging.Info("退票请求: ", req)
	tp := ticketpool.Tp
	for _, ticket := range req.Tickets {
		err := tp.RefundTickets(ticket.TrainId, ticket.StartStationId, ticket.DestStationId, ticket.StartTime, tp.GetIdBySeatTypeName(ticket.SeatType), fmt.Sprintf("%s %s", ticket.CarriageNumber, ticket.SeatNumber))
		if err != nil {
			logging.Error(err)
			return nil, err
		}
	}
	return &pb.RefundTicketResponse{
		IsOk: true,
	}, nil
}
