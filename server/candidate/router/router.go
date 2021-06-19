// Package router
// @Author LiuYong
// @Created at 2021-02-02
package router

import (
	v1 "candidate/controller/api/v1"
	"candidate/tools/setting"
	"common/middleware/token/usertoken"
	"common/router_tracer"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.Server.RunMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// token中间件
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
	apiV1 := r.Group("/candidate/api/v1")
	{
		apiV1.POST("/", v1.Candidate)
		apiV1.POST("/cash", v1.Cash)
		apiV1.POST("/state", v1.ReadState)
	}
	return r
}
