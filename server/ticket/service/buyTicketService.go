// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package service

import (
	"context"
	"google.golang.org/grpc"
	indentPb "rpc/indent/proto/indentRPC"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models/ticket"
)


func CheckUnHandleIndent(userId uint32) (bool, error) {
	indentConn, err := grpc.Dial("0.0.0.0:9440", grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer indentConn.Close()
	indentClient := indentPb.NewIndentServiceClient(indentConn)
	resp, err := indentClient.HasUnfinishedIndent(context.Background(), &indentPb.UnfinishedRequest{UserId: userId})
	if err != nil {
		return false, err
	}
	return resp.HasUnfinishedIndent, nil

}

func CheckConflict(passengerId *[]uint32 ,date string) (bool, error){
	isConflict, err := ticket.IsConflict(db, passengerId, date)
	if err != nil {
		return false, err
	}
	return isConflict, nil
}

func GetTickets(getTicketReq *ticketPoolPb.GetTicketRequest) ([]*ticketPoolPb.Ticket, error) {
	ticketPoolConn, err := grpc.Dial("0.0.0.0:9443", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer ticketPoolConn.Close()
	ticketPoolClient := ticketPoolPb.NewTicketPoolServiceClient(ticketPoolConn)

	tickets, err := ticketPoolClient.GetTicket(context.Background(), getTicketReq)
	if err != nil {
		return nil, err
	}
	return tickets.Tickets, nil
}

func CreateIndent(createIndentReq *indentPb.CreateRequest) (*indentPb.CreateResponse, error) {
	indentConn, err := grpc.Dial("0.0.0.0:9440", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer indentConn.Close()
	indentClient := indentPb.NewIndentServiceClient(indentConn)
	resp, err := indentClient.CreateIndent(context.Background(), createIndentReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

