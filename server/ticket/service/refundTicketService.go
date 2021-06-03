// @Author: KongLingWen
// @Created at 2021/2/17
// @Modified at 2021/2/17

package service

import (
	"common/tools/logging"
	"gorm.io/gorm"
	ticketPoolPb "rpc/ticketPool/proto/ticketPoolRPC"
	"ticket/models"
	"ticket/utils/database"
)

func RefundTicket(ticketsId []uint32) bool {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		res, err := models.DeleteTicketByTicketId(tx, ticketsId)
		if err != nil {
			logging.Error(err)
			return err
		}
		tickets := make([]*ticketPoolPb.Ticket, len(ticketsId))
		for i := 0; i < len(ticketsId); i++ {
			tickets[i] = &ticketPoolPb.Ticket{
				TrainId:        res[i].TrainId,
				StartStationId: res[i].StartStationId,
				StartTime:      res[i].StartTime.Format("2006-01-02"),
				DestStationId:  res[i].DestStationId,
				SeatType:       res[i].SeatType,
				CarriageNumber: res[i].CarriageNumber,
				SeatNumber:     res[i].SeatNumber,
			}
		}
		_, err = ts.tpCli.RefundTicket(&ticketPoolPb.RefundTicketRequest{Tickets: tickets})
		if err != nil {
			logging.Error(err)
			return err
		}
		return nil
	})

	if err != nil {
		logging.Error(err)
		return false
	}
	return true
}
