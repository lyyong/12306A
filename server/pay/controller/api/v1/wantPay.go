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

// 告知服务器想要通过支付宝支付 godoc
// @Summary 告知服务器想要支付
// @Description 获得支付需要的来自支付宝的OrderInfo
// @Accept json
// @Produce json
// @Param token query string true "认证信息"
// @Param wantPayR body v1.wantPayAbbRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{data=v1.wantPayAbbSend} "返回订单号和OrderInfo"
// @Failure 400 {object} controller.JSONResult{}
// @Router /want/abb [post]
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
	var payService service.PayService
	// 这里并没有接入支付宝所以只是用一个固定的OrderInfo
	wantPayS.OrderInfo = payService.WantPay(userInfo.UserId, wantPayR.OrderOutsideID)
	wantPayS.OrderOutsideID = wantPayR.OrderOutsideID
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, wantPayS))
}
