/*
* @Author: 余添能
* @Date:   2021/2/4 11:03 下午
 */
package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1:=r.Group("/search/v1/")
	{
		//查询所有站点
		v1.GET("/queryAllStations",)
		//查询两城市之间合适的车次及余票数量
		v1.GET("/queryTrainNos")
		//购票
		v1.POST("/buyTicket")
	}

	r.Run(":18081")
	return r
}

