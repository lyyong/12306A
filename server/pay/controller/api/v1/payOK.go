// Package v1
// @Author liuYong
// @Created at 2020-01-05
package v1

import (
	"common/middleware/token/usertoken"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
	"pay/service"
	"pay/tools/message"
)

type payOKAbbRecv struct {
	OrderInfo      string `json:"order_info" binding:"required"`
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

// PayOKAbb 使用支付宝支付完成通知服务器
func PayOKAbb(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var payOKR payOKAbbRecv
	if err := c.ShouldBindJSON(&payOKR); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	userInfo, ok := usertoken.GetUserInfo(c)
	if !ok {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	payService := service.NewPayService()
	err := payService.PayOK(userInfo.UserId, payOKR.OrderInfo, payOKR.OrderOutsideID)
	if err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	sender.Response(http.StatusOK, controller.NewJSONResult(message.PAYOK, noData))
}
