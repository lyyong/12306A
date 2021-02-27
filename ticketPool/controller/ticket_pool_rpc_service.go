// @Author LiuYong
// @Created at 2021-02-15
// @Modified at 2021-02-15
package controller

import (
	"context"
	"errors"
	"fmt"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"ticketPool/dao"
	"ticketPool/model/outer"
	"ticketPool/rdb"
	"ticketPool/rdb/service"
)

type TicketPoolRPCService struct {
}

func (t TicketPoolRPCService) GetTicket(ctx context.Context, request *ticketPoolRPC.GetTicketRequest) (*ticketPoolRPC.GetTicketResponse, error) {
	trainNo:=dao.GetTrainNumber(request.GetTrainId())
	date:=request.GetDate()
	startStation:=dao.GetStationName(request.GetStartStationId())
	endStation:=dao.GetStationName(request.GetDestStationId())
	passengers:=request.GetPassengers()
	var buyTickets []*outer.BuyTicket
	for _,passenger:=range passengers{
		seatType:=dao.GetSeatType(passenger.SeatTypeId)
		choose:=passenger.ChooseSeat
		buyTicket:=&outer.BuyTicket{}
		buyTicket.TrainNumber=trainNo
		buyTicket.StartStation=startStation
		buyTicket.EndStation=endStation
		buyTicket.SeatClass=seatType
		buyTicket.SeatPlace=choose
		buyTicket.StartTime=date
		buyTickets=append(buyTickets,buyTicket)
	}
	tickets:=service.BuyTicket(buyTickets)

	//返回结果
	response:=&ticketPoolRPC.GetTicketResponse{}
	if tickets==nil{
		response.Tickets=nil
		return response,errors.New("无票")
	}
	for i:=0;i<len(tickets);i++{
		ticket:=tickets[i]
		responseTicket:=&ticketPoolRPC.Ticket{}

		responseTicket.TrainId=request.TrainId
		responseTicket.TrainNum=trainNo
		responseTicket.StartStationId=request.StartStationId
		responseTicket.DestStationId=request.DestStationId
		responseTicket.StartTime=ticket.StartTime
		responseTicket.ArriveTime=ticket.EndTime
		responseTicket.StartStation=ticket.StartStation
		responseTicket.DestStation=ticket.EndStation

		responseTicket.PassengerId=passengers[i].PassengerId
		responseTicket.PassengerName=passengers[i].PassengerName
		responseTicket.SeatTypeId=passengers[i].SeatTypeId
		responseTicket.SeatType=dao.GetSeatType(passengers[i].SeatTypeId)
		responseTicket.SeatNumber=ticket.SeatNum
		responseTicket.CarriageNumber=ticket.CarriageNum
		responseTicket.Price=int32(ticket.Price)

		response.Tickets=append(response.Tickets,responseTicket)
	}
	return response,nil
	//request.TrainId
}

func (t TicketPoolRPCService) GetTicketNumber(ctx context.Context, request *ticketPoolRPC.GetTicketNumberRequest) (*ticketPoolRPC.GetTicketNumberResponse, error) {
	//panic("implement me")
	responses :=&ticketPoolRPC.GetTicketNumberResponse{}
	var trains []*ticketPoolRPC.TrainTicketInfo

	if responses==nil{
		fmt.Println("空请求")
		return nil, nil
	}
	date := request.Date
	conditions:=request.Condition
	for _,condition:=range conditions{
		//车次，上车站，下车站
		trainId:=condition.TrainId
		trainNo:=dao.GetTrainNumber(trainId)
		startStationId:=condition.StartStationId
		endStationId:=condition.DestStationId
		startStation:=dao.GetStationName(startStationId)
		endStation:=dao.GetStationName(endStationId)

		var seatInfos []*ticketPoolRPC.SeatInfo
		for i:=1;i<=len(dao.SeatTypes);i++{
			num:=service.QueryTicketNumByTrainNoAndDate(date,trainNo,dao.SeatTypes[uint32(i)],startStation,endStation)
			fmt.Println(num)

			seatInfo:=&ticketPoolRPC.SeatInfo{
				SeatTypeId: uint32(i),
				SeatNumber: int32(num),
			}
			seatInfos=append(seatInfos,seatInfo)
		}
		train:=&ticketPoolRPC.TrainTicketInfo{}
		train.TrainId=trainId
		train.SeatInfo=seatInfos
		trains=append(trains,train)
	}
	responses.TrainsTicketInfo=trains
	return responses,nil
}

func (t TicketPoolRPCService) RefundTicket(ctx context.Context, request *ticketPoolRPC.RefundTicketRequest) (*ticketPoolRPC.RefundTicketResponse, error) {
	response:=&ticketPoolRPC.RefundTicketResponse{}

	tickets:=request.Tickets
	if tickets==nil{
		return response,errors.New("无票可退")
	}
	for i:=0;i<len(request.Tickets);i++{
		ticket:=&outer.Ticket{}
		ticket.TrainNumber=tickets[i].TrainNum
		ticket.Date=tickets[i].StartTime

		ticket.StartTime=tickets[i].StartTime
		ticket.StartStation=tickets[i].StartStation
		ticket.EndStation=tickets[i].DestStation
		ticket.EndTime=tickets[i].ArriveTime

		ticket.CarriageNum=tickets[i].CarriageNumber
		ticket.SeatNum=tickets[i].SeatNumber
		ticket.SeatClass=tickets[i].SeatType

		trainMap,_:=rdb.RedisDB.HGetAll(ticket.TrainNumber).Result()
		ticket.StartStationNum=trainMap[ticket.StartStation]
		ticket.EndStationNum=trainMap[ticket.EndStation]

		success:=service.AddTicket(ticket)
		if success==false{
			response.IsOk=false
		}
	}
	if response.IsOk==false{
		return response,errors.New("部分成功")
	}
	return response,nil
}
