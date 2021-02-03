// @Author LiuYong
// @Created at 2021-02-03
// @Modified at 2021-02-03
package v1

import "github.com/gin-gonic/gin"

type candidateRecv struct {
	Date        string   `json:"date" binding:"required"`         // 发车日期 yyyy-mm-dd
	Time        string   `json:"time" binding:"required"`         // 发车时间 hh:mm
	TrainNumber string   `json:"train_number" binding:"required"` // 车次
	Passengers  []string `json:"passengers" binding:"required"`   // 乘客id
}

type candidateSend struct {
	OrderOutsideID string `json:"order_outside_id" binding:"required"` // 返回的订单编号
}

func Candidate(c *gin.Context) {

}
