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

// Reticket 告诉服务器执行退票请求 godoc
// @Summary 告诉服务器执行退票请求
// @Description 服务器执行退票请求
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
		return
	}
	// TODO 退票逻辑
	sender.Response(http.StatusOK, controller.NewJSONResult(message.OK, noDate))
}
