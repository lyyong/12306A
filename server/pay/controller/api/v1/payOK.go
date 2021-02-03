// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package v1

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/controller"
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
// @Param userID query string true "用户ID"
// @Param username query string true "用户名"
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
	// TODO 支付成功逻辑 查询支付宝看是否支付成功, 然后执行一些支付完成的逻辑
	sender.Response(http.StatusOK, controller.NewJSONResult(message.PAYOK, noData))
}
