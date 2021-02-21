// @Author: KongLingWen
// @Created at 2021/2/17
// @Modified at 2021/2/17

package controller

import (
	"common/tools/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)


type RefundRequest struct{
	TicketId 			int32	`json:"ticket_id"`
}

func RefundTicket(c *gin.Context){
	var refundReq RefundRequest
	if err := c.ShouldBindJSON(&refundReq); err != nil {
		logging.Error("bind param error:", err)
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: fmt.Sprintf("参数有误：%s", err.Error()), Data: nil})
		return
	}


}