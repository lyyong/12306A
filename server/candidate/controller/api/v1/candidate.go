// Package v1
// @Author LiuYong
// @Created at 2021-02-03
package v1

import (
	"candidate/controller"
	"candidate/service"
	"candidate/tools/message"
	"common/middleware/token/usertoken"
	"common/tools/logging"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type candidateRecv struct {
	Date           string `json:"date" binding:"required"`     // 发车日期 yyyy-mm-dd
	TrainID        uint   `json:"train_id" binding:"required"` // 车次
	StartStationID uint   `json:"start_station_id" binding:"required"`
	DestStationID  uint   `json:"dest_station_id" binding:"required"`
	ExpireDate     string `json:"expire_date" binding:"required"`
	SeatTypeID     int    `json:"seat_type_id" binding:"required"` // 座位类型, 0商务座, 1一等座, 2二等座
	Passengers     []uint `json:"passengers" binding:"required"`   // 乘客id
}

type candidateSend struct {
	OrderOutsideID string `json:"order_outside_id" binding:"required"` // 返回的订单编号
}

// Candidate 请求服务器执行候补功能
func Candidate(c *gin.Context) {
	send := controller.NewSend(c)
	noData := make(map[string]interface{})
	var cr candidateRecv
	if err := c.ShouldBindJSON(&cr); err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	// 验证接收参数
	if err := validateCandidateRecv(&cr); err != nil {
		logging.Info(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	userInfo, ok := usertoken.GetUserInfo(c)
	if !ok {
		// token数据出错
		send.Response(http.StatusOK, controller.NewJSONResult(message.PARAMS_ERROR, noData))
		return
	}
	// 创建service层
	canService, err := service.NewCandidateService()
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	// 缓存订单并且获得outsideID
	Date, err := time.Parse("2006-01-02", cr.Date)
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	ExpireDate, err := time.Parse("2006-01-02", cr.ExpireDate)
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	oosID, err := canService.CacheCandidate(userInfo.UserId, cr.TrainID, cr.StartStationID, cr.DestStationID, &Date, &ExpireDate, cr.Passengers, cr.SeatTypeID)
	if err != nil {
		logging.Error(err)
		send.Response(http.StatusOK, controller.NewJSONResult(message.ERROR, noData))
		return
	}
	var cs candidateSend
	cs.OrderOutsideID = oosID
	send.Response(http.StatusOK, controller.NewJSONResult(message.OK, cs))
}

// validateCandidateRecv 验证接收参数的正确性
func validateCandidateRecv(cr *candidateRecv) error {
	Date, err := time.Parse("2006-01-02", cr.Date)
	if err != nil {
		logging.Error(err)
		return err
	}
	ExpireDate, err := time.Parse("2006-01-02", cr.ExpireDate)
	if err != nil {
		logging.Error(err)
		return err
	}
	if !ExpireDate.After(time.Now().Add(24 * time.Hour)) {
		return errors.New("候补的截至时间不在1天后")
	}
	if !Date.After(ExpireDate.Add(23 * time.Hour)) {
		return errors.New("候补截至时间不在发车时间前1天")
	}
	if cr.DestStationID == cr.StartStationID {
		return errors.New("上下车站相同")
	}
	if len(cr.Passengers) <= 0 {
		return errors.New("没有乘客信息")
	}
	if cr.SeatTypeID < 0 || cr.SeatTypeID > 2 {
		return errors.New("座位类型出错")
	}
	return nil
}
