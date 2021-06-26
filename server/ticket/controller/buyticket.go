package controller

import (
	"common/middleware/token/usertoken"
	"common/tools/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	orderPb "rpc/pay/proto/orderRPCpb"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"rpc/user/userrpc"
	"ticket/service"
)

type BuyTicketRequest struct {
	TrainId        uint32      `json:"train_id"`
	StartStationId uint32      `json:"start_station_id"`
	DestStationId  uint32      `json:"dest_station_id"`
	Date           string      `json:"date"`
	Passengers     []Passenger `json:"passengers"`
}

type Passenger struct {
	PassengerId   uint32 `json:"passenger_id"`
	PassengerName string `json:"passenger_name"`
	SeatTypeId    uint32 `json:"seat_type_id"`
	ChooseSeat    string `json:"choose_seat"`
}

type BuyTicketResponse struct {
	OrderOuterId   string       `json:"order_id"`
	TrainId        uint32       `json:"train_id"`
	TrainNum       string       `json:"train_num"`
	StartStationId uint32       `json:"start_station_id"`
	StartStation   string       `json:"start_station"`
	StartTime      string       `json:"start_time"`
	DestStationId  uint32       `json:"dest_station_id"`
	DestStation    string       `json:"dest_station"`
	ArriveTime     string       `json:"arrive_time"`
	Date           string       `json:"date"`
	ExpiredTime    int32        `json:"expired_time"` // 单位秒
	Price          int32        `json:"price"`
	Tickets        []TicketInfo `json:"tickets"`
}

type TicketInfo struct {
	CertificateNumber string `json:"certificate_number"` // 乘客的身份证好
	PassengerName     string `json:"passenger_name"`
	SeatTypeId        uint32 `json:"seat_type_id"`
	SeatType          string `json:"seat_type"`
	CarriageNumber    string `json:"carriage_number"`
	SeatNumber        string `json:"seat_number"`
	Price             int32  `json:"price"`
}

func BuyTicket(c *gin.Context) {
	var btReq BuyTicketRequest
	if err := c.ShouldBindJSON(&btReq); err != nil {
		logging.Error("bind param error:", err)
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: fmt.Sprintf("参数有误：%s", err.Error()), Data: nil})
		return
	}
	// 获取用户信息
	userInfo, ok := usertoken.GetUserInfo(c)
	if !ok {
		logging.Error("token 错误")
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "token 错误", Data: nil})
		return
	}
	hasUnHandleIndent, err := service.CheckUnHandleIndent(uint32(userInfo.UserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 0, Msg: err.Error(), Data: nil})
		return
	}
	if hasUnHandleIndent {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "存在未处理订单", Data: nil})
		return
	}
	// 验证乘客身份信息
	allPassengerForUser, err := service.GetPassengers(uint32(userInfo.UserId))
	if err != nil {
		logging.Error(err)
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "验证乘客信息失败", Data: nil})
		return
	}
	allPassenger := make(map[uint32]*userrpc.Passenger)
	for i := range allPassengerForUser {
		allPassenger[uint32(allPassengerForUser[i].Id)] = allPassengerForUser[i]
	}
	passengerId := make([]uint32, len(btReq.Passengers))
	for index, value := range btReq.Passengers {
		_, isOK := allPassenger[value.PassengerId]
		if !isOK {
			c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "乘客信息有误", Data: nil})
			return
		}
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
	for index, value := range btReq.Passengers {
		passengers[index] = &ticketPoolPb.PassengerInfo{
			PassengerId:       value.PassengerId,
			PassengerName:     value.PassengerName,
			CertificateNumber: allPassenger[value.PassengerId].CertificateNumber,
			SeatTypeId:        value.SeatTypeId,
			ChooseSeat:        value.ChooseSeat,
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
	if err != nil || len(tickets) == 0 {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "出票失败", Data: nil})
		return
	}

	orderId := fmt.Sprintf("ticket_%d", uint32(userInfo.UserId))
	err = service.SaveTickets(orderId, tickets, 1800)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "缓存车票失败", Data: nil})
		return
	}

	createOrderReq := &orderPb.CreateRequest{
		UserID:         uint64(userInfo.UserId),
		Money:          8888,
		AffairID:       orderId,
		ExpireDuration: 1800,
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
			CertificateNumber: value.CertificateNumber,
			PassengerName:     value.PassengerName,
			SeatTypeId:        value.SeatTypeId,
			SeatType:          value.SeatType,
			CarriageNumber:    value.CarriageNumber,
			SeatNumber:        value.SeatNumber,
			Price:             value.Price,
		}
	}
	btResp := &BuyTicketResponse{
		OrderOuterId:   order.OrderOutsideID,
		TrainId:        btReq.TrainId,
		TrainNum:       tickets[0].TrainNum,
		StartStationId: btReq.StartStationId,
		StartStation:   tickets[0].StartStation,
		StartTime:      tickets[0].StartTime,
		DestStationId:  btReq.DestStationId,
		DestStation:    tickets[0].DestStation,
		ArriveTime:     tickets[0].ArriveTime,
		Date:           btReq.Date,
		ExpiredTime:    1800,
		Price:          8888,
		Tickets:        ticketsInfo,
	}
	c.JSON(http.StatusOK, Response{Code: 0, Msg: "出票成功", Data: btResp})
}
