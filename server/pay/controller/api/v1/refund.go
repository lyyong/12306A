// Package v1
// @Author LiuYong
// @Created at 2021-1-23
package v1

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
	"pay/tools/message"
)

type refundAbbRecv struct {
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

// Refund 退款
func Refund(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var refundR refundAbbRecv
	if err := c.ShouldBindJSON(&refundR); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	// TODO 通过outside id 查找到订单然后通过orderInfo向支付宝请求退款 过程可能是排队的
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, noData))
}
