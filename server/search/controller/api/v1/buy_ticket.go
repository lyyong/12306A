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
	buyTicket.Date=c.Query("date")
	buyTicket.TrainNo=c.Query("trainNo")
	buyTicket.StartStation=c.Query("startStation")
	buyTicket.EndStation=c.Query("endStation")
	buyTicket.SeatClass=c.Query("seatClass")
	//buyTicket.Date=c.PostForm("date")
	//buyTicket.TrainNo=c.PostForm("trainNo")
	//buyTicket.StartStation=c.PostForm("startStation")
	//buyTicket.EndStation=c.PostForm("endStation")
	//buyTicket.SeatClass=c.PostForm("seatClass")
	//fmt.Println(buyTicket)
	ticket:=rdb.BuyTicketByTrainNoAndDate(buyTicket)
	//fmt.Println(ticket)
	if ticket==nil{
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusNoContent,
			"ticket":ticket,
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"ticket":ticket,
		})
	}
}
