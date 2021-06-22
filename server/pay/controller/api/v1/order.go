// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package v1

import (
	"common/middleware/token/usertoken"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
	"pay/model"
	"pay/service"
	"pay/tools/message"
	"pay/tools/setting"
	"rpc/ticket/client"
	"rpc/ticket/proto/ticketRPC"
	"time"
)

type OrderInfo struct {
	TrainNumber    string `json:"train_number"`
	LeaveStation   string `json:"leave_station"`
	ArrivalStation string `json:"arrival_station"`
	LeaveTime      string `json:"leave_time"`
	TicketSum      int    `json:"ticket_sum"`
	FirstPassenger string `json:"first_passenger"`
}

type UnpayOrderInfo struct {
	OrderId        string            `json:"order_id"`
	TrainId        int               `json:"train_id"`
	TrainNum       string            `json:"train_num"`
	StartStationId int               `json:"start_station_id"`
	StartStation   string            `json:"start_station"`
	StartTime      string            `json:"start_time"`
	DestStationId  int               `json:"dest_station_id"`
	DestStation    string            `json:"dest_station"`
	ArriveTime     string            `json:"arrive_time"`
	Date           string            `json:"date"`
	ExpiredTime    int               `json:"expired_time"`
	Price          int               `json:"price"`
	Tickets        []UnpayTicketInfo `json:"tickets"`
}

type UnpayTicketInfo struct {
	PassengerId     int    `json:"passenger_id"`
	PassengerName   string `json:"passenger_name"`
	PassengerType   string `json:"passenger_type"`
	CertificateType string `json:"certificate_type"`
	SeatTypeId      int    `json:"seat_type_id"`
	SeatType        string `json:"seat_type"`
	CarriageNumber  string `json:"carriage_number"`
	SeatNumber      string `json:"seat_number"`
	Price           int    `json:"price"`
}

// GetUserHistoryOrders 获取用户的历史订单
func GetUserHistoryOrders(c *gin.Context) {
	userInfo, ok := usertoken.GetUserInfo(c)
	sender := controller.NewSend(c)
	if !ok {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, nil))
		return
	}

	// 获取历史订单
	orderService := service.NewOrderService()
	orders := orderService.GetOrdersByUserIDAndFinished(userInfo.UserId)
	getUserOrdersHelper(orders, sender)
}

// GetUserUnfinishedOrders 获取用户未完成的订单, 也就是还没有发车的订单
func GetUserUnfinishedOrders(c *gin.Context) {
	userInfo, ok := usertoken.GetUserInfo(c)
	sender := controller.NewSend(c)
	if !ok {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, nil))
		return
	}

	// 获取历史订单
	orderService := service.NewOrderService()
	orders := orderService.GetOrdersByUserIDAndUnfinished(userInfo.UserId)
	getUserOrdersHelper(orders, sender)
}

func getUserOrdersHelper(orders []*model.Order, sender *controller.Send) {
	// 没有数据
	if orders == nil || len(orders) == 0 {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, nil))
		return
	}
	// 请求获得相关票数据
	ticketCli, err := client.NewClientWithTarget(setting.RPCTarget.Ticket)
	if err != nil {
		logging.Error("创建连接ticket服务的rpc出错: ", err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, nil))
		return
	}
	orderOuterIDs := make([]string, len(orders))
	for i, a := range orders {
		orderOuterIDs[i] = a.OutsideID
	}
	resp, err := ticketCli.GetTicketByOrdersId(&ticketRPC.GetTicketByOrdersIdRequest{OrdersId: orderOuterIDs})
	if err != nil {
		logging.Error("获取关联票数据出错: ", err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, nil))
		return
	}

	orderInfo2Client := make([]*OrderInfo, 0, len(orders))
	for i := range orders {
		for j := range resp.GetList() {
			if len(resp.GetList()[j].Tickets) > 0 && resp.GetList()[j].Tickets[0].OrderOutsideId == orders[i].OutsideID {
				t := &OrderInfo{}
				t.TicketSum = len(resp.GetList()[j].Tickets)
				t.TrainNumber = resp.GetList()[j].Tickets[0].TrainNum
				t.LeaveStation = resp.GetList()[j].Tickets[0].StartStation
				t.ArrivalStation = resp.GetList()[j].Tickets[0].DestStation
				t.LeaveTime = resp.GetList()[j].Tickets[0].StartTime
				t.FirstPassenger = resp.GetList()[j].Tickets[0].PassengerName
				orderInfo2Client = append(orderInfo2Client, t)
			}
		}
	}

	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, orderInfo2Client))
}

// GetUserUnpayOrders 获取用户未支付的订单
func GetUserUnpayOrders(c *gin.Context) {
	userInfo, ok := usertoken.GetUserInfo(c)
	sender := controller.NewSend(c)
	if !ok {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, nil))
		return
	}

	// 获取历史订单
	orderService := service.NewOrderService()
	order := orderService.GetOrdersByUserIDAndUnpay(userInfo.UserId)
	// 没有数据
	if order == nil {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, nil))
		return
	}
	// 请求获得相关票数据
	ticketCli, err := client.NewClientWithTarget(setting.RPCTarget.Ticket)
	if err != nil {
		logging.Error("创建连接ticket服务的rpc出错: ", err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, nil))
		return
	}
	resp, err := ticketCli.GetUnHandleTickets(&ticketRPC.GetUnHandleTicketsRequest{
		UserId: uint32(order.UserID),
	})
	if err != nil {
		logging.Error("获取关联票数据出错: ", err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, nil))
		return
	}

	if resp == nil || len(resp.Tickets) == 0 {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, nil))
		return
	}
	// TODO 数据补全
	orderInfo2Client := &UnpayOrderInfo{
		OrderId:        order.OutsideID,
		TrainId:        0,
		TrainNum:       resp.Tickets[0].TrainNum,
		StartStationId: 0,
		StartStation:   resp.Tickets[0].StartStation,
		StartTime:      resp.Tickets[0].StartTime,
		DestStationId:  0,
		DestStation:    resp.Tickets[0].DestStation,
		ArriveTime:     resp.Tickets[0].DestTime,
		Date:           resp.Tickets[0].StartTime,
		ExpiredTime:    1800 - order.CreatedAt.Add(time.Minute*30).Second() + order.CreatedAt.Second(),
		Price:          order.Money,
	}

	for i := range resp.Tickets {
		orderInfo2Client.Tickets = append(orderInfo2Client.Tickets, UnpayTicketInfo{
			PassengerName:  resp.Tickets[i].PassengerName,
			SeatType:       resp.Tickets[i].SeatType,
			CarriageNumber: resp.Tickets[i].CarriageNumber,
			SeatNumber:     resp.Tickets[i].SeatNumber,
			Price:          int(resp.Tickets[i].Price),
		})
	}

	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, orderInfo2Client))
}

// CancelUnpayOrder 取消用户未支付的订单
func CancelUnpayOrder(c *gin.Context) {
	userInfo, ok := usertoken.GetUserInfo(c)
	sender := controller.NewSend(c)
	if !ok {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, nil))
		return
	}

	orderServer := service.NewOrderService()
	orderServer.CancelUnpayOrder(userInfo.UserId)
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, "取消成功"))
}
