// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package service

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"rpc/pay/client/orderRPCClient"
	orderPb "rpc/pay/proto/orderRPCpb"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models"
	"ticket/utils/redispool"
	"time"
)

func CheckConflict(passengerId *[]uint32 ,date string) (bool, error){
	isConflict, err := models.IsConflict(passengerId, date)
	if err != nil {
		return false, err
	}
	return isConflict, nil
}

func GetTickets(getTicketReq *ticketPoolPb.GetTicketRequest) ([]*ticketPoolPb.Ticket, error) {
	tickets, err := ts.tpCli.GetTicket(getTicketReq)
	if err != nil {
		return nil, err
	}
	return tickets.Tickets, nil
}

func CheckUnHandleIndent(userId uint32) (bool, error) {
	resp, err := ts.orderCli.GetNoFinishOrder(&orderPb.SearchCondition{UserID: uint64(userId)})
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
	resp, err := ts.orderCli.Create(createReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SaveTickets(userId uint32, tickets []*ticketPoolPb.Ticket, expireTime int32) error {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	data, err := json.Marshal(tickets)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%d_ticket", userId)
	conn.Do("SET", key, data, "EX", expireTime)
	return nil
}

func payOK(payOKInfo *orderRPCClient.PayOKOrderInfo) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("%d_ticket", payOKInfo.UserID)

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		// 通知支付模块
		return
	}
	var tpTickets []*ticketPoolPb.Ticket
	err = json.Unmarshal(data, &tpTickets)
	if err != nil {
		// 通知支付模块
		return
	}
	tickets := make([]models.Ticket, len(tpTickets))
	for i := 0; i < len(tpTickets); i++ {
		startTime, _ := time.Parse("2006-01-02 15:04", tpTickets[i].StartTime)
		arriveTime, _ := time.Parse("2006-01-02 15:04", tpTickets[i].ArriveTime)

		tickets[i] = models.Ticket{
			Model:          gorm.Model{},
			UserId:         uint32(payOKInfo.UserID),
			TrainId:        tpTickets[i].TrainId,
			TrainNum:       tpTickets[i].TrainNum,
			StartStationId: tpTickets[i].StartStationId,
			StartStation:   tpTickets[i].StartStation,
			StartTime:      startTime,
			DestStationId:  tpTickets[i].DestStationId,
			DestStation:    tpTickets[i].DestStation,
			DestTime:       arriveTime,
			SeatType:       tpTickets[i].SeatType,
			CarriageNumber: tpTickets[i].CarriageNumber,
			SeatNumber:     tpTickets[i].SeatNumber,
			Price:          tpTickets[i].Price,
			OrderOutsideId: payOKInfo.OutsideID,
			PassengerName:  tpTickets[i].PassengerName,
			PassengerId:    tpTickets[i].PassengerId,
			State:          0,
		}
	}
	err = models.AddMultipleTicket(&tickets)
	if err != nil {
		// 通知支付模块
		return
	}
}