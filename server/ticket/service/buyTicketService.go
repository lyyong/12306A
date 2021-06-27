// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package service

import (
	"common/tools/logging"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"rpc/pay/client/orderRPCClient"
	orderPb "rpc/pay/proto/orderRPCpb"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"rpc/user/userrpc"
	"ticket/models"
	"ticket/utils/redispool"
	"time"
)

func GetPassengers(userID uint32) ([]*userrpc.Passenger, error) {
	return ts.userCli.ListPassenger(uint(userID))
}

func CheckConflict(passengerId *[]uint32, date string) (bool, error) {
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
	resp, err := ts.orderCli.ExistNoFinishOrder(&orderPb.SearchCondition{UserID: uint64(userId)})
	if err != nil {
		return false, err
	}
	return resp.Exist, err
}

func CreateOrder(createReq *orderPb.CreateRequest) (*orderPb.CreateRespond, error) {
	resp, err := ts.orderCli.Create(createReq)
	if err != nil {
		logging.Error(err)
		return nil, err
	}
	return resp, nil
}

func SaveTickets(key string, tickets []*ticketPoolPb.Ticket, expireTime int32) error {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	data, err := json.Marshal(tickets)
	if err != nil {
		logging.Error(err)
		return err
	}

	conn.Do("SET", key, data, "EX", expireTime)
	return nil
}

func payOK(payOKInfo *orderRPCClient.PayOKOrderInfo) {
	conn := redispool.RedisPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("ticket_%d", payOKInfo.UserID)

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logging.Error(err)
		return
	}
	var tpTickets []*ticketPoolPb.Ticket
	err = json.Unmarshal(data, &tpTickets)
	if err != nil {
		logging.Error(err)
		return
	}
	tickets := make([]*models.Ticket, len(tpTickets))
	for i := 0; i < len(tpTickets); i++ {
		startTime, _ := time.ParseInLocation("2006-01-02 15:04", tpTickets[i].StartTime, time.Local)
		arriveTime, _ := time.ParseInLocation("2006-01-02 15:04", tpTickets[i].ArriveTime, time.Local)

		tickets[i] = &models.Ticket{
			Model:             gorm.Model{},
			UserId:            uint32(payOKInfo.UserID),
			TrainId:           tpTickets[i].TrainId,
			TrainNum:          tpTickets[i].TrainNum,
			StartStationId:    tpTickets[i].StartStationId,
			StartStation:      tpTickets[i].StartStation,
			StartTime:         startTime,
			DestStationId:     tpTickets[i].DestStationId,
			DestStation:       tpTickets[i].DestStation,
			DestTime:          arriveTime,
			SeatType:          tpTickets[i].SeatType,
			CarriageNumber:    tpTickets[i].CarriageNumber,
			SeatNumber:        tpTickets[i].SeatNumber,
			Price:             tpTickets[i].Price,
			OrderOutsideId:    payOKInfo.OutsideID,
			PassengerId:       tpTickets[i].PassengerId,
			PassengerName:     tpTickets[i].PassengerName,
			CertificateNumber: tpTickets[i].CertificateNumber,
			State:             0,
		}
	}
	err = models.AddMultipleTicket(tickets)
	if err != nil {
		logging.Error(err)
		return
	}
}
