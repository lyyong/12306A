/*
* @Author: 余添能
* @Date:   2021/2/4 11:03 下午
 */
package router

import (
	v12 "12306A-search/controller/api/v1"
	"12306A-search/tools/settings"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// r := gin.Default()
	gin.SetMode(settings.Server.RunMode)
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/search/api/v1/")
	{
		// 查询所有站点
		v1.GET("/stations", v12.QueryAllStation)
		// 查询车次的所有站点
		v1.GET("/station", v12.QueryStationByTrainNo)
		// 查询两城市之间合适的车次及余票数量
		v1.GET("/remainder", v12.QueryRemainder)

		v1.GET("/remainder/:train_id/:date/:start_station_id/:end_station_id", v12.QueryRemainderWithTrainNumber)
	}

	r.Run(":18081")

	return r
}
