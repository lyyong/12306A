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
	//request.TrainId
}

func (t TicketPoolRPCService) GetTicketNumber(ctx context.Context, request *ticketPoolRPC.GetTicketNumberRequest) (*ticketPoolRPC.GetTicketNumberResponse, error) {
	//panic("implement me")
	//trainIds := request.TrainId
	//date := request.Date
	//startStationId := request.StartStationId
	//destStationId := request.DestStationId
	//startStation:=dao.GetStationName(startStationId)
	//endStation:=dao.GetStationName(destStationId)
	//
	//responses :=&ticketPoolRPC.GetTicketNumberResponse{}
	//var trains []*ticketPoolRPC.TrainTicketInfo
	//
	//for _,trainId:=range trainIds{
	//	trainNumber:=dao.GetTrainNumber(trainId)
	//	//获取余票
	//	businessSeatNum:=service.QueryTicketNumByTrainNoAndDate(date,trainNumber,"businessSeat",startStation,endStation)
	//	firstSeatNum:=service.QueryTicketNumByTrainNoAndDate(date,trainNumber,"businessSeat",startStation,endStation)
	//	secondSeatNum:=service.QueryTicketNumByTrainNoAndDate(date,trainNumber,"businessSeat",startStation,endStation)
	//	//只有高铁座位
	//	var seatInfos[]*ticketPoolRPC.SeatInfo
	//	bSeatInfo:=&ticketPoolRPC.SeatInfo{}
	//	bSeatInfo.SeatNumber=int32(businessSeatNum)
	//	bSeatInfo.SeatTypeId=1
	//	seatInfos=append(seatInfos,bSeatInfo)
	//
	//	fSeatInfo:=&ticketPoolRPC.SeatInfo{}
	//	fSeatInfo.SeatNumber=int32(firstSeatNum)
	//	fSeatInfo.SeatTypeId=2
	//	seatInfos=append(seatInfos,fSeatInfo)
	//
	//	sSeatInfo:=&ticketPoolRPC.SeatInfo{}
	//	sSeatInfo.SeatNumber=int32(secondSeatNum)
	//	sSeatInfo.SeatTypeId=3
	//	seatInfos=append(seatInfos,sSeatInfo)
	//
	//	train:=&ticketPoolRPC.TrainTicketInfo{}
	//	train.TrainId=trainId
	//	train.SeatInfo=seatInfos
	//	trains=append(trains,train)
	//}
	//responses.TrainsTicketInfo=trains
	return nil,nil
}

func (t TicketPoolRPCService) RefundTicket(ctx context.Context, request *ticketPoolRPC.RefundTicketRequest) (*ticketPoolRPC.RefundTicketResponse, error) {
	panic("implement me")
}
