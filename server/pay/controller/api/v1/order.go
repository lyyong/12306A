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
)

type OrderInfo struct {
	TrainNumber    string `json:"train_number"`
	LeaveStation   string `json:"leave_station"`
	ArrivalStation string `json:"arrival_station"`
	LeaveTime      string `json:"leave_time"`
	TicketSum      int    `json:"ticket_sum"`
	FirstPassenger string `json:"first_passenger"`
}

// @Summary 用户获取自己的历史订单信息
// @Description
// @Accept json
// @Produce json
// @Param token header string true "认证信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /history [get]
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

// @Summary 用户获取自己未出行的订单信息
// @Description
// @Accept json
// @Produce json
// @Param token header string true "认证信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /unfished [get]
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

// @Summary 用户获取自己的未支付订单
// @Description
// @Accept json
// @Produce json
// @Param token header string true "认证信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /unpay [get]
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

	if len(resp.Tickets) == 0 || len(resp.Tickets) == 0 {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, nil))
		return
	}

	orderInfo2Client := &OrderInfo{}
	ticketSum := len(resp.Tickets)
	orderInfo2Client.TrainNumber = resp.Tickets[0].TrainNum
	orderInfo2Client.LeaveStation = resp.Tickets[0].StartStation
	orderInfo2Client.ArrivalStation = resp.Tickets[0].DestStation
	orderInfo2Client.LeaveTime = resp.Tickets[0].StartTime
	orderInfo2Client.FirstPassenger = resp.Tickets[0].PassengerName
	orderInfo2Client.TicketSum = ticketSum

	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, orderInfo2Client))
}
