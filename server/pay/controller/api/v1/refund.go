// @Author LiuYong
// @Created at 2021-1-23
// @Modified at 2021-1-23
package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
	"pay/tools/logging"
	"pay/tools/message"
)

type refundRecv struct {
	Username       string `json:"username" binding:"required"`
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

func Refund(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var refundR refundRecv
	if err := c.ShouldBindJSON(&refundR); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, message.PARAMS_ERROR, noData)
	}
	// TODO 通过outside id 查找到订单然后通过orderInfo向支付宝请求退款 过程可能是排队的
	sender.Response(http.StatusOK, message.OK, noData)
}
