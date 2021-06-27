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
	OrderId        string       `json:"order_id"`
	TrainId        int          `json:"train_id"`
	TrainNum       string       `json:"train_num"`
	StartStationId int          `json:"start_station_id"`
	StartStation   string       `json:"start_station"`
	StartTime      string       `json:"start_time"`
	DestStationId  int          `json:"dest_station_id"`
	DestStation    string       `json:"dest_station"`
	ArriveTime     string       `json:"arrive_time"`
	Date           string       `json:"date"`
	ExpiredTime    int          `json:"expired_time"`
	Price          int          `json:"price"`
	Tickets        []TicketInfo `json:"tickets"`
}

type TicketInfo struct {
	TicketID          int    `json:"ticket_id"`
	CertificateNumber string `json:"certificate_number"`
	PassengerName     string `json:"passenger_name"`
	SeatType          string `json:"seat_type"`
	CarriageNumber    string `json:"carriage_number"`
	SeatNumber        string `json:"seat_number"`
	Price             int    `json:"price"`
}

const (
	ticketPaySuccessful = iota
	ticketFinish
	ticketRefund
	ticketChange
	ticketWaitCash
	ticketChanged
)

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

	orderInfos2Client := make([]*OrderInfo, 0, len(orders))
	for i := range orders {
		for j := range resp.GetList() {
			if len(resp.GetList()[j].Tickets) > 0 && resp.GetList()[j].Tickets[0].OrderOutsideId == orders[i].OutsideID {
				tickets := make([]TicketInfo, 0, len(resp.List[j].Tickets))
				for k := range resp.List[j].Tickets {
					if resp.List[j].Tickets[k].State == ticketChange || resp.List[j].Tickets[k].State == ticketRefund || resp.List[j].Tickets[k].State == ticketWaitCash {
						continue
					}
					tickets = append(tickets, TicketInfo{
						TicketID:          int(resp.List[j].Tickets[k].Id),
						CertificateNumber: resp.List[j].Tickets[k].CertificateNumber,
						PassengerName:     resp.List[j].Tickets[k].PassengerName,
						SeatType:          resp.List[j].Tickets[k].SeatType,
						CarriageNumber:    resp.List[j].Tickets[k].CarriageNumber,
						SeatNumber:        resp.List[j].Tickets[k].SeatNumber,
						Price:             int(resp.List[j].Tickets[k].Price),
					})
				}
				if len(tickets) == 0 {
					continue
				}
				t := &OrderInfo{
					OrderId:        resp.List[j].Tickets[0].OrderOutsideId,
					TrainId:        int(resp.List[j].Tickets[0].TrainId),
					TrainNum:       resp.List[j].Tickets[0].TrainNum,
					StartStationId: int(resp.List[j].Tickets[0].StartStationId),
					StartStation:   resp.List[j].Tickets[0].StartStation,
					StartTime:      resp.List[j].Tickets[0].StartTime,
					DestStationId:  int(resp.List[j].Tickets[0].DestStationId),
					DestStation:    resp.List[j].Tickets[0].DestStation,
					ArriveTime:     resp.List[j].Tickets[0].DestTime,
					Date:           resp.List[j].Tickets[0].StartTime,
					Price:          int(resp.List[j].Tickets[0].Price),
					Tickets:        tickets,
				}

				orderInfos2Client = append(orderInfos2Client, t)
			}
		}
	}

	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, orderInfos2Client))
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
	orderInfo2Client := &OrderInfo{
		OrderId:        order.OutsideID,
		TrainId:        int(resp.Tickets[0].TrainId),
		TrainNum:       resp.Tickets[0].TrainNum,
		StartStationId: int(resp.Tickets[0].StartStationId),
		StartStation:   resp.Tickets[0].StartStation,
		StartTime:      resp.Tickets[0].StartTime,
		DestStationId:  int(resp.Tickets[0].DestStationId),
		DestStation:    resp.Tickets[0].DestStation,
		ArriveTime:     resp.Tickets[0].DestTime,
		Date:           resp.Tickets[0].StartTime,
		ExpiredTime:    1800 - order.CreatedAt.Add(time.Minute*30).Second() + order.CreatedAt.Second(),
		Price:          order.Money,
	}

	for i := range resp.Tickets {
		orderInfo2Client.Tickets = append(orderInfo2Client.Tickets, TicketInfo{
			TicketID:          int(resp.Tickets[i].Id),
			CertificateNumber: resp.Tickets[i].CertificateNumber,
			PassengerName:     resp.Tickets[i].PassengerName,
			SeatType:          resp.Tickets[i].SeatType,
			CarriageNumber:    resp.Tickets[i].CarriageNumber,
			SeatNumber:        resp.Tickets[i].SeatNumber,
			Price:             int(resp.Tickets[i].Price),
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
