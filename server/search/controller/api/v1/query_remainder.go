/*
* @Author: 余添能
* @Date:   2021/2/4 10:56 下午
 */
package v1

import (
	"12306A-search/rdb"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": message.PARAMS_ERROR,
			"msg":  message.GetMsg(message.PARAMS_ERROR),
		})
		return
	}
	//fmt.Println(date,startCity,endCity,"aaa")
	//fmt.Println(search)
	//trains:=rdb.Query(search)
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
