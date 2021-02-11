package controller

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticket/service"
)

type testStruct struct {
	Username string `json:"user_name"`
	UserId int `json:"user"`
}

func BuyTicket(c *gin.Context){
/*  测试数据（JSON）
{
	"user_id": 0,
	"train_id": 1,
	"start_station_id": 2,
	"dest_station_id" : 5,
	"date" : "2021-02-05",
	"passengers" : [
		{
			"passenger_id" : 12,
			"seat_type_id" : 0,
			"choose_seat" : "A"
		},
		{
			"passenger_id" : 13,
			"seat_type_id" : 0,
			"choose_seat" : "B"
		},
		{
			"passenger_id" : 14,
			"seat_type_id" : 1,
			"choose_seat" : "C"
		}
	]
}
*/
	var btReq service.BuyTicketRequest

	if err := c.ShouldBindJSON(&btReq); err != nil {
		// 参数有误
		logging.Error("bind param error:", err)
		c.JSON(http.StatusBadRequest, Response{
			Code: 0,
			Msg:  "参数有误",
			Data: err,
		})
		return
	}

	btResp, err := service.BuyTicket(&btReq)

	if err != nil || btResp == nil{
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  "出票失败",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "出票成功",
		Data: btResp,
	})
}