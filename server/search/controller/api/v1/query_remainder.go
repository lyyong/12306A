/*
* @Author: 余添能
* @Date:   2021/2/4 10:56 下午
 */
package v1

import (
	"12306A-search/rdb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pay/tools/message"
	"time"
)

type QueryRemainderReq struct {
	StartCity string `json:"start_city"`
	EndCity   string `json:"end_city"`
	Date      string `json:"date"`
	Type      string `json:"type"`
}

type QueryRemainderWithTrainNumReq struct {
	TrainID        uint32    `uri:"train_id" binding:"required"`
	Date           time.Time `uri:"date" binding:"required" time_format:"2006-01-02"`
	StartStationID uint32    `uri:"start_station_id" binding:"required"`
	EndStationID   uint32    `uri:"end_station_id" binding:"required"`
}

func QueryRemainder(c *gin.Context) {
	var req QueryRemainderReq
	req.StartCity = c.Query("start_city")
	req.EndCity = c.Query("end_city")
	req.Date = c.Query("date")
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": message.PARAMS_ERROR,
	//		"msg":  message.GetMsg(message.PARAMS_ERROR),
	//	})
	//	fmt.Println("abc,err")
	//	return
	//}
	//fmt.Println("aaa",req)
	trains := rdb.QueryTicketNumByDate(req.Date, req.StartCity, req.EndCity)
	if trains == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": message.OK,
			"msg":  "车次票数",
			"data": gin.H{"list": "没有车次"},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "车次票数",
			"data": gin.H{"list": trains},
		})
	}

}

func QueryRemainderWithTrainNumber(c *gin.Context) {
	var req QueryRemainderWithTrainNumReq

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": message.PARAMS_ERROR,
			"msg":  message.GetMsg(message.PARAMS_ERROR),
		})
		fmt.Println("参数解析错误")
		//logging.Info("参数解析错误")
		return
	}

	train := rdb.QueryTicketNumByDateWithTrainNumber(req.TrainID, req.StartStationID,
		req.EndStationID, req.Date.Format("2006-01-02"))
	c.JSON(http.StatusOK, gin.H{
		"code": message.OK,
		"msg":  message.GetMsg(message.OK),
		"data": train,
	})

}
