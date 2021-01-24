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

type wantPayRecv struct {
	Username       string `json:"username" binding:"required"`
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

type wantPaySend struct {
	OrderOutsideID string `json:"order_outside_id"`
	OrderInfo      string `json:"order_info"`
}

func WantPay(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var wantPayR wantPayRecv
	if err := c.ShouldBindJSON(&wantPayR); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, message.PARAMS_ERROR, noData)
	}
	var wantPayS wantPaySend
	// TODO 亲求支付宝获取OrderInfo, 然后填充wantPayS
	sender.Response(http.StatusOK, message.OK, wantPayS)
}
