// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
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

// 告知服务器支付宝支付完成了 godoc
// @Summary 告知服务器通过支付宝支付完成
// @Description 服务器对支付进行验证
// @Accept json
// @Produce json
// @Param token header string true "认证信息"
// @Param wantPayR body v1.payOKAbbRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /ok/abb [post]
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
