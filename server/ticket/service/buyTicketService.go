// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package service

import (
	"common/tools/logging"
	"context"
	"google.golang.org/grpc"
	indentPb "rpc/indent/proto/indentRPC"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
)

type BuyTicketRequest struct{
	UserId 			int32	`json:"user_id"`  	// 暂时先由前端传参，后续通过用户身份认证获得
	TrainId			int32 	`json:"train_id"`
	StartStationId	int32 	`json:"start_station_id"`
	DestStationId	int32 	`json:"dest_station_id"`
	Date			string 	`json:"date"`
	Passengers  	[]Passenger `json:"passengers"`
}

type Passenger struct {
	PassengerId int32 	`json:"passenger_id"`
	SeatTypeId	int32 	`json:"seat_type_id"`
	ChooseSeat  string 	`json:"choose_seat"`
}

type BuyTicketResponse struct {
	IndentOuterId 	string `json:"indent_id"`
	TrainId 		int32 `json:"train_id"`
	StartStationId	int32 `json:"start_station_id"`
	StartTime 		string `json:"start_time"`
	DestStationId	int32 `json:"dest_station_id"`
	ArriveTime 		string `json:"arrive_time"`
	Date			string `json:"date"`
	ExpiredTime 	int32 `json:"expired_time"` // second
	Amount 			int32 `json:"amount"`
	Tickets			[]TicketInfo `json:"tickets"`
}

type TicketInfo struct {
	PassengerId int32 `json:"passenger_id"`
	SeatTypeId 	int32 `json:"seat_type_id"`
	CarriageNumber string `json:"carriage_number"`
	SeatNumber string `json:"seat_number"`
	Amount int32 `json:"amount"`
}

func BuyTicket(btReq *BuyTicketRequest) (*BuyTicketResponse, error) {
	logging.Info("调用订单模块")

	// 调用订单模块 -- 验证用户（user）是否有未处理的订单
	indentConn, err := grpc.Dial("0.0.0.0:9440", grpc.WithInsecure())
	if err != nil {
		logging.Error("连接订单模块失败", err)
		return nil, err
	}
	defer indentConn.Close()
	indentClient := indentPb.NewIndentServiceClient(indentConn)
	resp, err := indentClient.HasUnfinishedIndent(context.Background(), &indentPb.UnfinishedRequest{
		UserId: btReq.UserId,
	})
	if resp.HasUnfinishedIndent || err != nil {
		return nil, err
	}

	logging.Info("查询是否冲突")
	// 根据 passenger_id 查询 Ticket 表 -- 验证每个乘车人（Passenger) 是否已购买时间冲突的车票
	passengerId := make([]int32, len(btReq.Passengers))
	for index, value := range btReq.Passengers{
		passengerId[index] = value.PassengerId
	}
	isConflict, err := IsConflict(&passengerId)
	if isConflict || err != nil {
		return nil, err
	}

	logging.Info("调用票池出票")
	// 调用票池 -- 出票
	ticketPoolConn, err := grpc.Dial("0.0.0.0:9443", grpc.WithInsecure())
	if err != nil {
		logging.Error("连接票池失败", err)
		return nil, err
	}
	defer ticketPoolConn.Close()
	ticketPoolClient := ticketPoolPb.NewTicketPoolServiceClient(ticketPoolConn)

	passengers := make([]*ticketPoolPb.PassengerInfo, len(btReq.Passengers))
	for index, value := range btReq.Passengers{
		passengers[index] = &ticketPoolPb.PassengerInfo{
			PassengerId: value.PassengerId,
			SeatTypeId:  value.SeatTypeId,
			ChooseSeat:  value.ChooseSeat,
		}
	}

	tickets, err := ticketPoolClient.GetTicket(context.Background(), &ticketPoolPb.GetTicketRequest{
		TrainId:        btReq.TrainId,
		StartStationId: btReq.StartStationId,
		DestStationId:  btReq.DestStationId,
		Date:           btReq.Date,
		Passengers:     passengers,
	})
	if err != nil {
		// 票池出票失败
		return nil, err
	}

	logging.Info("调用订单模块创建订单")
	// 调用订单模块 -- 创建订单
	createIndentResp, err := indentClient.CreateIndent(context.Background(), &indentPb.CreateRequest{
		UserId:         btReq.UserId,
		TrainId:        btReq.TrainId,
		StartStationId: btReq.StartStationId,
		StartTime:      tickets.Tickets[0].StartTime,
		DestStationId:  btReq.DestStationId,
		ArriveTime:     tickets.Tickets[0].ArriveTime,
		Date:           tickets.Tickets[0].Date,
		ExpiredTime:	1800,
		TicketNumber:   int32(len(tickets.Tickets)),
		Amount:         8850,
	})
	if err != nil {
		// 创建订单失败

		return nil, err
	}

	// 订单模块创建订单成功，以 userId 为 key 将 tickets.Tickets 存入 redis， 过期时间为 订单过期时间 + rpc调用延迟时间

	logging.Info("出票成功，返回数据")
	// 返回
	ticketsInfo := make([]TicketInfo, len(tickets.Tickets))
	for index, value := range tickets.Tickets {
		ticketsInfo[index] = TicketInfo{
			PassengerId:    value.PassengerId,
			SeatTypeId:     value.SeatTypeId,
			CarriageNumber: value.CarriageNumber,
			SeatNumber:     value.SeatNumber,
			Amount:         value.Amount,
		}
	}
	btResp := &BuyTicketResponse{
		IndentOuterId:  createIndentResp.IndentOuterId,
		TrainId:        btReq.TrainId,
		StartStationId: btReq.StartStationId,
		StartTime:      tickets.Tickets[0].StartTime,
		DestStationId:  btReq.DestStationId,
		ArriveTime:     tickets.Tickets[0].ArriveTime,
		Date:           btReq.Date,
		ExpiredTime:    1800,
		Amount:         8850,
		Tickets:        ticketsInfo,
	}
	return btResp, nil
}
