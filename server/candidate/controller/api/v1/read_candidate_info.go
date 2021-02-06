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

type readRecv struct {
	CandidateID string `json:"candidate_id"` // 候补id
}

type readSend struct {
	CandidateID string `json:"candidate_id"`
	Passengers  []struct {
		Name              string `json:"name"`               // 乘客名
		CertificateNumber string `json:"certificate_number"` // 证件号
		Type              int    `json:"type"`               // 乘客类型 0为普通乘客, 1为学生
	} `json:"passengers"`
	state int // 状态 0为正在候补,1为候补成功,2为候补失败
}

// ReadState 请求服务器查看候补状态 godoc
// @Summary 请求服务器查看候补状态
// @Description 查看自己的候补订单
// @Accept json
// @Produce json
// @Param userID query string true "用户ID"
// @Param username query string true "用户名"
// @Param wantPayR body v1.readRecv true "需要接受的信息"
// @Success 200 {object} controller.JSONResult{data=v1.readSend} "返回成功"
// @Failure 400 {object} controller.JSONResult{}
// @Router /state [post]
func ReadState(c *gin.Context) {
	send := controller.NewSend(c)
	noData := make(map[string]interface{})
	var rr readRecv
	if err := c.ShouldBindJSON(&rr); err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	var rs readSend
	// TODO 获取信息
	send.Response(http.StatusOK, controller.NewJSONResult(message.OK, rs))
}
