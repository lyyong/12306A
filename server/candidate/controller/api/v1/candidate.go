// @Author LiuYong
// @Created at 2021-02-03
// @Modified at 2021-02-03
package v1

import (
	"candidate/controller"
	"candidate/service"
	"candidate/tools/message"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type candidateRecv struct {
	Date       string   `json:"date" binding:"required"`       // 发车日期 yyyy-mm-dd
	TrainID    int      `json:"train_id" binding:"required"`   // 车次
	Passengers []string `json:"passengers" binding:"required"` // 乘客id
}

type candidateSend struct {
	OrderOutsideID string `json:"order_outside_id" binding:"required"` // 返回的订单编号
}

// Candidate 请求服务器执行候补功能 godoc
// @Summary 请求服务器执行候补功能
// @Description 发送需要候补的信息给服务器, 服务器将执行候补功能
// @Accept json
// @Produce json
// @Param userID query int true "用户ID"
// @Param username query string true "用户名"
// @Param wantPayR body v1.candidateRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{data=v1.candidateSend} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router / [post]
func Candidate(c *gin.Context) {
	send := controller.NewSend(c)
	noData := make(map[string]interface{})
	var cr candidateRecv
	if err := c.ShouldBindJSON(&cr); err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	canService, err := service.NewCandidateService()
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	oosID, err := canService.CacheCandidate(userID, cr.TrainID, cr.Date, cr.Passengers)
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	var cs candidateSend
	cs.OrderOutsideID = oosID
	send.Response(http.StatusOK, controller.NewJSONResult(message.OK, cs))
}
