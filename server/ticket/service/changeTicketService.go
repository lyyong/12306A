// Package service
// @Author LiuYong
// @Created at 2021-06-23
package service

import (
	"common/tools/logging"
	"rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models"
	"ticket/utils/database"
	"time"
)

func Change(userID, ticketID, trainID, startStationID, destStationID uint, date string) *models.Ticket {
	// 检查票是否属于用户
	oldTicket, err := models.GetTicketByID(ticketID)
	if err != nil {
		logging.Error(err)
		return nil
	}
	if oldTicket.UserId != uint32(userID) || oldTicket.TrainId != uint32(trainID) || oldTicket.State != models.TicketBuySuccessful {
		logging.Error("信息验证失败")
		return nil
	}
	newDate, err := time.Parse("2006-01-02", date)
	if err != nil || newDate.Before(time.Now()) {
		logging.Error("时间匹配出错")
		return nil
	}

	// 检查站点是否在同一个城市
	if startStationID != uint(oldTicket.StartStationId) {
		var oldStartCity, newStartCity string
		database.DB.Raw("select city from stations where id = ?", startStationID).Scan(&newStartCity)
		database.DB.Raw("select city from stations where id = ?", oldTicket.StartStationId).Scan(&oldStartCity)
		if newStartCity != oldStartCity {
			logging.Error("城市匹配出错")
			return nil
		}
	}
	if destStationID != uint(oldTicket.DestStationId) {
		var oldDestCity, newDestCity string
		database.DB.Raw("select city from stations where id = ?", destStationID).Scan(&newDestCity)
		database.DB.Raw("select city from stations where id = ?", oldTicket.DestStationId).Scan(&oldDestCity)
		if newDestCity != oldDestCity {
			logging.Error("城市匹配出错")
			return nil
		}
	}
	// 获取票
	getTicketReq := &ticketPoolRPC.GetTicketRequest{
		TrainId:        uint32(trainID),
		StartStationId: uint32(startStationID),
		DestStationId:  uint32(destStationID),
		Date:           date,
		Passengers: []*ticketPoolRPC.PassengerInfo{&ticketPoolRPC.PassengerInfo{
			PassengerId:   oldTicket.PassengerId,
			PassengerName: oldTicket.PassengerName,
			SeatTypeId:    getSeatTypeId(oldTicket.SeatType),
			ChooseSeat:    oldTicket.SeatNumber[:1],
		}},
	}
	getTicketResp, err := ts.tpCli.GetTicket(getTicketReq)
	if err != nil || len(getTicketResp.Tickets) == 1 {
		return nil
	}

	// 退票
	refundTicketReq := &ticketPoolRPC.RefundTicketRequest{
		Tickets: []*ticketPoolRPC.Ticket{&ticketPoolRPC.Ticket{
			Id:             uint32(oldTicket.ID),
			TrainId:        oldTicket.TrainId,
			TrainNum:       oldTicket.TrainNum,
			StartStationId: oldTicket.StartStationId,
			StartStation:   oldTicket.StartStation,
			StartTime:      oldTicket.StartTime.Format("2006-01-02"),
			DestStationId:  oldTicket.DestStationId,
			DestStation:    oldTicket.DestStation,
			SeatTypeId:     getSeatTypeId(oldTicket.SeatType),
			SeatType:       oldTicket.SeatType,
			CarriageNumber: oldTicket.CarriageNumber,
			SeatNumber:     oldTicket.SeatNumber,
			PassengerName:  oldTicket.PassengerName,
			PassengerId:    oldTicket.PassengerId,
			OrderId:        oldTicket.OrderOutsideId,
			Price:          oldTicket.Price,
		}},
	}
	refundResp, err := ts.tpCli.RefundTicket(refundTicketReq)
	if err != nil || !refundResp.IsOk {
		// TODO 退票失败, 退会新拿的票

		return nil
	}
	startTime, _ := time.Parse("2006-01-02 15:04", getTicketResp.Tickets[0].StartTime)
	arriveTime, _ := time.Parse("2006-01-02 15:04", getTicketResp.Tickets[0].ArriveTime)
	newTicket := models.Ticket{
		UserId:         uint32(userID),
		TrainId:        uint32(trainID),
		TrainNum:       getTicketResp.Tickets[0].TrainNum,
		StartStationId: uint32(startStationID),
		StartStation:   getTicketResp.Tickets[0].StartStation,
		StartTime:      startTime,
		DestStationId:  uint32(destStationID),
		DestStation:    getTicketResp.Tickets[0].DestStation,
		DestTime:       arriveTime,
		SeatType:       getTicketResp.Tickets[0].SeatType,
		CarriageNumber: getTicketResp.Tickets[0].CarriageNumber,
		SeatNumber:     getTicketResp.Tickets[0].SeatNumber,
		Price:          getTicketResp.Tickets[0].Price,
		OrderOutsideId: oldTicket.OrderOutsideId,
		PassengerName:  oldTicket.PassengerName,
		PassengerId:    oldTicket.PassengerId,
		State:          models.TicketRefundFinish,
	}
	_, err = models.DeleteTicketByTicketId(database.DB, []uint32{uint32(oldTicket.ID)})
	if err != nil {
		logging.Error(err)
		return nil
	}
	err = models.AddMultipleTicket(&[]models.Ticket{newTicket})
	if err != nil {
		logging.Error(err)
		return nil
	}
	return &newTicket
}

func getSeatTypeId(seatType string) uint32 {
	switch seatType {
	case "商务座":
		return 0
	case "一等座":
		return 1
	case "二等座":
		return 2
	default:
		return 3
	}
}
