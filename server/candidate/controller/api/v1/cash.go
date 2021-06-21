// Package v1
// @Author LiuYong
// @Created at 2021-02-03
package v1

import (
	"candidate/controller"
	"candidate/service"
	"candidate/tools/message"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 兑现候补

type cashRecv struct {
	CandidateID []uint `json:"candidate_id" binding:"required"`
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
	cc, err := service.NewCandidateService()
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	err = cc.Cash(cr.CandidateID)
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	send.Response(http.StatusOK, controller.NewJSONResult(message.OK, noData))
}
