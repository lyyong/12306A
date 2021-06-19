// Package v1
// @Author LiuYong
// @Created at 2021-02-03
package v1

import (
	"candidate/controller"
	"candidate/tools/message"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 兑现候补

type cashRecv struct {
	CandidateID string `json:"candidate_id" binding:"required"`
}

// Cash 请求服务器兑现候补
func Cash(c *gin.Context) {
	send := controller.NewSend(c)
	noData := make(map[string]interface{})
	var cr cashRecv
	if err := c.ShouldBindJSON(&cr); err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	// TODO 兑现逻辑
	send.Response(http.StatusOK, controller.NewJSONResult(message.OK, noData))
}
