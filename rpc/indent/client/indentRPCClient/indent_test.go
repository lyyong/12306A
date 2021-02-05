package indentRPCClient
// @Author KongLingWen
// @Created at 2021/1/30
// @Modified at 2021/1/30

//package indentRPCClient
//
//import (
//	"context"
//	"fmt"
//	"google.golang.org/grpc"
//	pb "rpc2/indent/proto/indentRPC"
//	"testing"
//)
//
//func TestCreate(t *testing.T) {
//	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
//	if err != nil {
//		fmt.Println("连接服务器失败", err)
//	}
//	defer conn.Close()
//
//	c := pb.NewIndentServiceClient(conn)
//	ticket1 := pb.Ticket{
//		TrainId:        0,
//		StartStation:   "苏州",
//		DestStation:    "上海",
//		Date:           "2021-01-29",
//		SeatType:       "硬座",
//		CarriageNumber: "10",
//		SeatNumber:     "01",
//		PassengerId:    "0",
//	}
//	ticket2 := pb.Ticket{
//		TrainId:        0,
//		StartStation:   "苏州",
//		DestStation:    "上海",
//		Date:           "2021-01-29",
//		SeatType:       "硬座",
//		CarriageNumber: "10",
//		SeatNumber:     "02",
//		PassengerId:    "1",
//	}
//	tickets := []*pb.Ticket{&ticket1, &ticket2}
//
//	cr := pb.CreateRequest{
//		UserId:       1,
//		TicketNumber: 2,
//		Tickets: tickets,
//	}
//	createResponse, err := c.Create(context.Background(), &cr)
//	if err != nil {
//		fmt.Println("Can not create indent:", err)
//		return
//	}
//	fmt.Println("get Response")
//	fmt.Println("[Client]response:",createResponse)
//}
