// @Author LiuYong
// @Created at 2021-1-23
// @Modified at 2021-1-23
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

type wantPayAbbRecv struct {
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

type wantPayAbbSend struct {
	OrderOutsideID string `json:"order_outside_id"`
	OrderInfo      string `json:"order_info"`
}

// WantPayAbb 告知服务器想要使用支付宝支付, 然后获取支付密钥
func WantPayAbb(c *gin.Context) {
	sender := controller.NewSend(c)
	noData := make(map[string]interface{})
	var wantPayR wantPayAbbRecv
	if err := c.ShouldBindJSON(&wantPayR); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	// 检查该用户是持有该订单
	userInfo, ok := usertoken.GetUserInfo(c)
	if !ok {
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}

	var wantPayS wantPayAbbSend
	payService := service.NewPayService()
	// 这里并没有接入支付宝所以只是用一个固定的OrderInfo
	wantPayS.OrderInfo = payService.WantPay(userInfo.UserId, wantPayR.OrderOutsideID)
	wantPayS.OrderOutsideID = wantPayR.OrderOutsideID
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, wantPayS))
}
