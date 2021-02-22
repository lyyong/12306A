// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	orderPb "rpc/pay/proto/orderRPCpb"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models/ticket"
)

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

func CheckUnHandleIndent(userId uint32) (bool, error) {
	orderConn, err := grpc.Dial("0.0.0.0:8082", grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer orderConn.Close()
	orderClient := orderPb.NewOrderRPCServiceClient(orderConn)
	resp, err := orderClient.GetNoFinishOrder(context.Background(), &orderPb.SearchCondition{UserID: uint64(userId)})
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, nil
	}else {
		return true, nil
	}
}

func CreateOrder(createReq *orderPb.CreateRequest) (*orderPb.CreateRespond, error) {
	orderConn, err := grpc.Dial("0.0.0.0:8082", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer orderConn.Close()
	orderClient := orderPb.NewOrderRPCServiceClient(orderConn)
	resp, err := orderClient.Create(context.Background(), createReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SaveTickets(userId uint32, tickets []*ticketPoolPb.Ticket, expireTime int32) error {
	conn := redisPool.Get()
	defer conn.Close()
	data, err := json.Marshal(tickets)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%d_ticket", userId)
	conn.Do("SET", key, data, "EX", expireTime)
	return nil
}