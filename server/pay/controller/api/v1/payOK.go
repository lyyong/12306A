package v1

import (
	"12306A/server/pay/controller"
	"12306A/server/pay/tools/logging"
	"12306A/server/pay/tools/message"
	"github.com/gin-gonic/gin"
	"net/http"
)

type payOKRecv struct {
	Username       string `json:"username" binding:"required"`
	OrderInfo      string `json:"order_info" binding:"required"`
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

func PayOK(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var payOKr payOKRecv
	if err := c.ShouldBindJSON(&payOKr); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, message.ERROR, noData)
	}
	// TODO 支付成功逻辑
	sender.Response(http.StatusOK, message.PAYOK, noData)
}
