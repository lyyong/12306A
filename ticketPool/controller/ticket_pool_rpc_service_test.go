/*
* @Author: 余添能
* @Date:   2021/2/28 12:00 上午
 */
package controller

import (
	"context"
	"fmt"
	pb "rpc/ticketPool/proto/ticketPoolRPC"
	"testing"
	"ticketPool/dao"
)


func TestTicketPoolRPCService_GetTicketNumber(t *testing.T) {
	//rpc.Setup()
	ticket:=&TicketPoolRPCService{}
	var conditions []*pb.GetTicketNumberRequest_Condition
	condition:=&pb.GetTicketNumberRequest_Condition{}
	condition.TrainId=dao.GetTrainId("G21")
	condition.StartStationId=dao.GetStationId("北京南")
	condition.DestStationId=dao.GetStationId("上海虹桥")
	conditions=append(conditions,condition)
	request:=&pb.GetTicketNumberRequest{
		Date: "2021-02-25",
		Condition:conditions,
	}
	response,_:=ticket.GetTicketNumber(context.Background(),request)
	//fmt.Println(ticket.GetTicketNumber(context.Background(),request))
	if response==nil{
		fmt.Println("无票")
	}
	fmt.Println(response)
	//TicketPoolRPCService_GetTicketNumber()
}

func TestTicketPoolRPCService_GetTicket(t *testing.T) {
	ticket:=&TicketPoolRPCService{}
	request:=&pb.GetTicketRequest{
		Date: "2021-02-25",
		TrainId: dao.GetTrainId("G21"),
		StartStationId: dao.GetStationId("北京南"),
		DestStationId: dao.GetStationId("上海虹桥"),
	}
	var passagers []*pb.PassengerInfo
	p1:=&pb.PassengerInfo{
		SeatTypeId: 0,
	}
	passagers=append(passagers,p1)
	request.Passengers=passagers
	fmt.Println(ticket.GetTicket(context.Background(),request))
}

func TestTicketPoolRPCService_RefundTicket(t *testing.T) {
	ticketService:=&TicketPoolRPCService{}
	request:=&pb.RefundTicketRequest{}
	var tickets []*pb.Ticket
	ticket:=&pb.Ticket{}

	ticket.TrainId=dao.GetTrainId("G21")
	ticket.TrainNum="G21"
	ticket.StartStationId=dao.GetStationId("北京南")
	ticket.StartStation="北京南"
	ticket.DestStationId=dao.GetStationId("上海虹桥")
	ticket.DestStation="上海虹桥"
	ticket.SeatTypeId=0
	ticket.SeatType="businessSeat"
	ticket.StartTime="2021-02-25 16:00:00"
	ticket.ArriveTime="2021-02-25 21:14:00"
	ticket.CarriageNumber="29"
	ticket.SeatNumber="2A"
	tickets=append(tickets,ticket)
	request.Tickets=tickets
	ticketService.RefundTicket(context.Background(),request)

}