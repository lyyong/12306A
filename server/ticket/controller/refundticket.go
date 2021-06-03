// @Author: KongLingWen
// @Created at 2021/2/17
// @Modified at 2021/2/17

package controller

import (
	"common/tools/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticket/service"
)


type RefundRequest struct{
	TicketsId 			[]uint32	`json:"tickets_id"`
}

func RefundTicket(c *gin.Context){
	var refundReq RefundRequest
	if err := c.ShouldBindJSON(&refundReq); err != nil {
		logging.Error("bind param error:", err)
		c.JSON(http.StatusBadRequest, Response{Code: 0, Msg: fmt.Sprintf("请求参数有误：%s", err.Error()), Data: nil})
		return
	}

	isOk := service.RefundTicket(refundReq.TicketsId)
	if !isOk {
		c.JSON(http.StatusInternalServerError, Response{Code: 0, Msg: "退票失败", Data: nil})
	}
	c.JSON(http.StatusOK, Response{Code: 0, Msg: "退票成功", Data: nil})
}