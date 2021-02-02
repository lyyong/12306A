// @Author LiuYong
// @Created at 2021-02-02
// @Modified at 2021-02-02
package v1

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"reticket/controller"
	"reticket/tools/message"
)

type reticketRecv struct {
	TicketOutsideID string `json:"ticket_outside_id" binding:"required"`
}

// 告知服务器支付宝支付完成了 godoc
// @Summary 告知服务器通过支付宝支付完成
// @Description 服务器对支付进行验证
// @Accept json
// @Produce json
// @Param userID query string true "用户ID"
// @Param username query string true "用户名"
// @Param wantPayR body v1.reticketRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router / [post]
func Reticket(c *gin.Context) {
	sender := controller.NewSend(c)
	noDate := make(map[string]interface{})
	var rr reticketRecv
	if err := c.ShouldBindJSON(&rr); err != nil {
		logging.Error(err)
		sender.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noDate))
	}
	// TODO 退票逻辑
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, noDate))
}
