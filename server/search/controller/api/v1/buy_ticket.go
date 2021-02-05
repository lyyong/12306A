/*
* @Author: 余添能
* @Date:   2021/2/4 10:56 下午
 */
package v1

import (
	"12306A/server/search/model/outer"
	"12306A/server/search/rdb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BuyTicket(c *gin.Context)  {
	buyTicket:=&outer.BuyTicket{}
	buyTicket.Date=c.PostForm("date")
	buyTicket.TrainNo=c.PostForm("trainNo")
	buyTicket.StartStation=c.PostForm("startStation")
	buyTicket.EndStation=c.PostForm("endStation")
	buyTicket.SeatClass=c.PostForm("seatClass")
	ticket:=rdb.BuyTicketByTrainNoAndDate(buyTicket)
	//fmt.Println(ticket)
	c.JSON(http.StatusOK,gin.H{"ticket":ticket})
}
