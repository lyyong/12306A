// Package v1
// @Author LiuYong
// @Created at 2021-02-03
package v1

import (
	"candidate/controller"
	"candidate/model"
	"candidate/service"
	"candidate/tools/message"
	"common/middleware/token/usertoken"
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

type readSend struct {
	OrderID      string      `json:"order_id"`
	TrainNumber  string      `json:"train_number"`
	StartStation string      `json:"start_station"`
	DestStation  string      `json:"dest_station"`
	Passengers   []passenger `json:"passengers"`
	State        int         `json:"state"` // 与model.Candidate的状态相同
}

type passenger struct {
	ID          uint   `json:"id"`
	CandidateID uint   `json:"candidate_id"`
	Name        string `json:"name"` // 乘客名
	Type        int    `json:"type"` // 乘客类型 0为普通乘客, 1为学生
}

// ReadState 请求服务器查看候补状态
func ReadState(c *gin.Context) {
	send := controller.NewSend(c)
	noData := make(map[string]interface{})
	// 查看用户的token信息
	userInfo, ok := usertoken.GetUserInfo(c)
	if !ok {
		logging.Info("请求header中解析token出错")
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
	}
	cc, _ := service.NewCandidateService()
	cans := cc.ReadCandidate(userInfo.UserId)
	if len(cans) == 0 {
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	passengers := make([]passenger, len(cans))
	for i := range passengers {
		passengers[i] = passenger{
			ID:          cans[i].PassengerID,
			CandidateID: cans[i].ID,
			Name:        cans[i].PassengerName,
			Type:        0,
		}
	}
	rs := readSend{
		OrderID:      cans[0].OrderID,
		TrainNumber:  model.GetTrainNumber(cans[0].TrainID),
		StartStation: model.GetStationName(cans[0].StartStationID),
		DestStation:  model.GetStationName(cans[0].DestStationID),
		Passengers:   passengers,
		State:        cans[0].State,
	}
	send.Response(http.StatusOK, controller.NewJSONResult(message.OK, rs))
}
