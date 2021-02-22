package controller

import (
	"common/tools/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	orderPb "rpc/pay/proto/orderRPCpb"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/service"
)

type BuyTicketRequest struct{
	UserId 			uint32	`json:"user_id"`  	// 暂时先由前端传参，后续通过用户身份认证获得
	TrainId			uint32 	`json:"train_id"`
	StartStationId	uint32 	`json:"start_station_id"`
	DestStationId	uint32 	`json:"dest_station_id"`
	Date			string 	`json:"date"`
	Passengers  	[]Passenger `json:"passengers"`
}

type Passenger struct {
	PassengerId uint32 	`json:"passenger_id"`
	SeatTypeId	uint32 	`json:"seat_type_id"`
	ChooseSeat  string 	`json:"choose_seat"`
}

type BuyTicketResponse struct {
	IndentOuterId 	string `json:"indent_id"`
	TrainId 		uint32 `json:"train_id"`
	StartStationId	uint32 `json:"start_station_id"`
	StartTime 		string `json:"start_time"`
	DestStationId	uint32 `json:"dest_station_id"`
	ArriveTime 		string `json:"arrive_time"`
	Date			string `json:"date"`
	ExpiredTime 	int32 `json:"expired_time"` // 单位秒
	Amount 			int32 `json:"amount"`
	Tickets			[]TicketInfo `json:"tickets"`
}

type TicketInfo struct {
	PassengerId uint32 `json:"passenger_id"`
	SeatTypeId 	uint32 `json:"seat_type_id"`
	CarriageNumber string `json:"carriage_number"`
	SeatNumber string `json:"seat_number"`
	Amount int32 `json:"amount"`
}

func BuyTicket(c *gin.Context){
	var btReq BuyTicketRequest
	if err := c.ShouldBindJSON(&btReq); err != nil {
		logging.Error("bind param error:", err)
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: fmt.Sprintf("参数有误：%s", err.Error()), Data: nil})
		return
	}

	hasUnHandleIndent, err := service.CheckUnHandleIndent(btReq.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 0, Msg: err.Error(), Data: nil})
		return
	}
	if hasUnHandleIndent {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "存在未处理订单", Data: nil})
		return
	}

	passengerId := make([]uint32, len(btReq.Passengers))
	for index, value := range btReq.Passengers{
		passengerId[index] = value.PassengerId
	}
	isConflict, err := service.CheckConflict(&passengerId, btReq.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 0, Msg: err.Error(), Data: nil})
		return
	}
	if isConflict {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "与已购车票冲突", Data: nil})
		return
	}

	passengers := make([]*ticketPoolPb.PassengerInfo, len(btReq.Passengers))
	for index, value := range btReq.Passengers{
		passengers[index] = &ticketPoolPb.PassengerInfo{
			PassengerId: value.PassengerId,
			SeatTypeId:  value.SeatTypeId,
			ChooseSeat:  value.ChooseSeat,
		}
	}
	getTicketReq := &ticketPoolPb.GetTicketRequest{
		TrainId:        btReq.TrainId,
		StartStationId: btReq.StartStationId,
		DestStationId:  btReq.DestStationId,
		Date:           btReq.Date,
		Passengers:     passengers,
	}
	tickets, err := service.GetTickets(getTicketReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "出票失败", Data: nil})
		return
	}

	err = service.SaveTickets(btReq.UserId, tickets, 1800)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "缓存车票失败", Data: nil})
		return
	}

	createOrderReq := &orderPb.CreateRequest{
		UserID:         uint64(btReq.UserId),
		Money:          8888,
		AffairID:       "",
		ExpireDuration: 0,
		CreatedBy:      "ticket",
	}
	order, err := service.CreateOrder(createOrderReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "创建订单失败", Data: nil})
		return
	}

	ticketsInfo := make([]TicketInfo, len(tickets))
	for index, value := range tickets {
		ticketsInfo[index] = TicketInfo{
			PassengerId:    value.PassengerId,
			SeatTypeId:     value.SeatTypeId,
			CarriageNumber: value.CarriageNumber,
			SeatNumber:     value.SeatNumber,
			Amount:         value.Amount,
		}
	}
	btResp := &BuyTicketResponse{
		IndentOuterId:  order.OrderOutsideID,
		TrainId:        btReq.TrainId,
		StartStationId: btReq.StartStationId,
		StartTime:      tickets[0].StartTime,
		DestStationId:  btReq.DestStationId,
		ArriveTime:     tickets[0].ArriveTime,
		Date:           btReq.Date,
		ExpiredTime:    1800,
		Amount:         8850,
		Tickets:        ticketsInfo,
	}
	c.JSON(http.StatusOK, Response{Code: 0, Msg: "出票成功", Data: btResp})
}