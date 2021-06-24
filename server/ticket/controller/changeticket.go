// Package controller
// @Author LiuYong
// @Created at 2021-06-23
package controller

import (
	"common/middleware/token/usertoken"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticket/service"
)

type changeTicketRecv struct {
	TicketId       uint   `json:"ticket_id" binding:"required"` // 改签的票id
	TrainId        uint   `json:"train_id" binding:"required"`  // 改签的车次
	StartStationId uint   `json:"start_station_id" binding:"required"`
	DestStationId  uint   `json:"dest_station_id" binding:"required"`
	Date           string `json:"date" binding:"required"` // 改签后的日期
}

type changeTicketSend struct {
	TicketId      uint   `json:"ticket_id"`
	PassengerName string `json:"passenger_name"`
	TrainNumber   string `json:"train_number"`
	StartStation  string `json:"start_station"`
	DestStation   string `json:"dest_station"`
	StartTime     string `json:"start_time"`
}

// ChangeTicket 改签票
func ChangeTicket(c *gin.Context) {
	userInfo, ok := usertoken.GetUserInfo(c)
	if !ok {
		logging.Info("token获取用户信息出错")
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "出票失败", Data: nil})
		return
	}
	recv := &changeTicketRecv{}
	if err := c.ShouldBindJSON(recv); err != nil {
		logging.Info("参数错误")
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "出票失败", Data: nil})
		return
	}
	newTicket := service.Change(userInfo.UserId, recv.TicketId, recv.TrainId, recv.StartStationId, recv.DestStationId, recv.Date)
	if newTicket == nil {
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: "出票失败", Data: nil})
		return
	}
	send := &changeTicketSend{
		TicketId:      newTicket.ID,
		PassengerName: newTicket.PassengerName,
		TrainNumber:   newTicket.TrainNum,
		StartStation:  newTicket.StartStation,
		DestStation:   newTicket.DestStation,
		StartTime:     newTicket.StartTime.Format("2006-01-02;15:04:05"),
	}
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "改签成功",
		Data: send,
	})
}
