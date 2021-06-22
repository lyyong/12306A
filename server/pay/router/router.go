// Package router
// @Author liuYong
// @Created at 2020-01-05
package router

import (
	"common/middleware/token/usertoken"
	"common/router_tracer"
	v1 "pay/controller/api/v1"
	"pay/tools/setting"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化一个路由
func InitRouter() *gin.Engine {
	gin.SetMode(setting.Server.RunMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(usertoken.TokenParser())
	// 设置使用链路追踪
	if router_tracer.IsTracing() {
		r.Use(func(context *gin.Context) {
			cli, _ := router_tracer.GetClient()
			spin := (*cli.Tracer()).StartSpan(context.FullPath())
			defer spin.Finish()
			context.Next()
		})
	}
	payApiV1 := r.Group("/pay/api/v1")
	{
		okGroup := payApiV1.Group("/ok")
		{
			okGroup.POST("/abb", v1.PayOKAbb)
		}
		wantGroup := payApiV1.Group("/want")
		{
			wantGroup.POST("/abb", v1.WantPayAbb)
		}
	}
	orderApiV1 := r.Group("/order/api/v1")
	{
		orderApiV1.GET("/history", v1.GetUserHistoryOrders)
		orderApiV1.GET("/unpay", v1.GetUserUnpayOrders)
		orderApiV1.GET("/unfinished", v1.GetUserUnfinishedOrders)
		orderApiV1.DELETE("/unpay", v1.CancelUnpayOrder)
	}

	return r
}
