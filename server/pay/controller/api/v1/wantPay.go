// @Author LiuYong
// @Created at 2021-1-23
// @Modified at 2021-1-23
package v1

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
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
// @Param userID query string true "用户ID"
// @Param username query string true "用户名"
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
	var wantPayS wantPayAbbSend
	// TODO 亲求支付宝获取OrderInfo, 然后填充wantPayS
	wantPayS.OrderInfo = "asdiuyUYGFYGV7567hgvfhjv"
	wantPayS.OrderOutsideID = wantPayR.OrderOutsideID
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, wantPayS))
}
