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

type refundAbbRecv struct {
	OrderOutsideID string `json:"order_outside_id" binding:"required"`
}

// 退款 godoc
// @Summary 请求退款
// @Description 给订单号然后进行退款操作
// @Accept json
// @Produce json
// @Param token query string true "认证信息"
// @Param refundR body v1.refundAbbRecv true "订单信息"
// @Success 200 {object} controller.JSONResult{} "成功信息"
// @Failure 400 {object} controller.JSONResult{}
// @Router /refund/abb [post]
func RefundAbb(c *gin.Context) {
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
