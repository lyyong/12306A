/*
* @Author: 余添能
* @Date:   2021/2/4 11:03 下午
 */
package router

import (
	v12 "12306A/server/search/controller/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1:=r.Group("/search/v1/")
	{
		//查询所有站点
		v1.GET("/queryAllStations",v12.QueryAllStation)
		//查询两城市之间合适的车次及余票数量
		v1.GET("/queryTrainNos",v12.Query)
		//购票
		v1.POST("/buyTicket",v12.BuyTicket)
	}

	r.Run(":18081")
	return r
}

