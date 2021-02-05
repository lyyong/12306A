// @Author LiuYong
// @Created at 2021-02-03
// @Modified at 2021-02-03
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

// Candidate 请求服务器执行候补功能 godoc
// @Summary 请求服务器执行候补功能
// @Description 发送需要候补的信息给服务器, 服务器将执行候补功能
// @Accept json
// @Produce json
// @Param userID query string true "用户ID"
// @Param username query string true "用户名"
// @Param wantPayR body v1.cashRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /cash [post]
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
