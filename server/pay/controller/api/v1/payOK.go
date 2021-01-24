// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
	"pay/tools/logging"
	"pay/tools/message"
)

type payOKRecv struct {
	Username       string `json:"username" binding:"required"`
	OrderInfo      string `json:"order_info" binding:"required"`
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

// @Summary 支付
func PayOK(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var payOKR payOKRecv
	if err := c.ShouldBindJSON(&payOKR); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, message.ERROR, noData)
	}
	// TODO 支付成功逻辑 查询支付宝看是否支付成功, 然后执行一些支付完成的逻辑
	sender.Response(http.StatusOK, message.PAYOK, noData)
}
