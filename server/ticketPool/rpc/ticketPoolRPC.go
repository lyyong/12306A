package rpc

import (
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
	tp := ticketpool.Tp
	train := tp.GetTrain(req.TrainId)

	seatCountMap := make(map[uint32]int32)
	for i := 0; i < len(req.Passengers); i++ {
		seatTypeId := req.Passengers[i].SeatTypeId
		seatCountMap[seatTypeId]++
	}

	seatsMap, err := tp.GetTicket(req.TrainId, req.StartStationId, req.DestStationId, req.Date, seatCountMap)
	if err != nil {
		return &pb.GetTicketResponse{Tickets: nil}, err
	}

	startStation := train.GetStopStation(req.StartStationId)
	destStation := train.GetStopStation(req.DestStationId)
	st, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", req.Date, startStation.StartTime))
	at, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", req.Date, destStation.ArriveTime))
	if at.Before(st) {
		at.AddDate(0, 0, 1)
	}
	startTime := st.Format("2006-01-02 15:04")
	arriveTime := at.Format("2006-01-02 15:04")

	tickets := make([]*pb.Ticket, len(req.Passengers))
	ticketIndex := 0
	for seatTypeId, seats := range seatsMap {
		seatIndex := 0
		seatType := tp.GetSeatInfo(seatTypeId).SeatType
		for i := 0; i < len(req.Passengers); i++ {
			passengerInfo := req.Passengers[i]
			if seatTypeId == passengerInfo.SeatTypeId {
				carriageAndSeat := strings.Split(seats[seatIndex], " ")
				// 生成车票信息

				tickets[ticketIndex] = &pb.Ticket{
					Id:             0,
					TrainId:        req.TrainId,
					TrainNum:       train.TrainNum,
					StartStationId: req.StartStationId,
					StartStation:   startStation.StationName,
					StartTime:      startTime,
					DestStationId:  req.DestStationId,
					DestStation:    destStation.StationName,
					ArriveTime:     arriveTime,
					SeatTypeId:     seatTypeId,
					SeatType:       seatType,
					CarriageNumber: carriageAndSeat[0],
					SeatNumber:     carriageAndSeat[1],
					PassengerName:  passengerInfo.PassengerName,
					PassengerId:    passengerInfo.PassengerId,
					OrderId:        0,
					Price:          88,
				}
				ticketIndex++
				seatIndex++
			}
		}
	}
	return &pb.GetTicketResponse{Tickets: tickets}, nil
}

func (tps *TicketPoolServer) GetTicketNumber(ctx context.Context, req *pb.GetTicketNumberRequest) (*pb.GetTicketNumberResponse, error) {
	// 查询车次余票
	tp := ticketpool.Tp
	trainsId := req.TrainId
	tti := make([]*pb.TrainTicketInfo, len(trainsId))

	for i := 0; i < len(trainsId); i++ {
		seatCountMap := tp.SearchTicketCount(trainsId[i], req.StartStationId, req.DestStationId, req.Date)
		seatInfo := make([]*pb.SeatInfo, len(seatCountMap))
		index := 0
		for seatTypeId, count := range seatCountMap {
			seatInfo[index] = &pb.SeatInfo{SeatTypeId: seatTypeId, SeatNumber: count}
			index++
		}
		tti[i] = &pb.TrainTicketInfo{TrainId: trainsId[i], SeatInfo: seatInfo}
	}

	return &pb.GetTicketNumberResponse{TrainsTicketInfo: tti},nil
}

func (tps *TicketPoolServer) RefundTicket(ctx context.Context, req *pb.RefundTicketRequest) (*pb.RefundTicketResponse, error) {
	// 退票到票池

	return &pb.RefundTicketResponse{
		IsOk: true,
	}, nil
}