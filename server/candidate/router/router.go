// @Author LiuYong
// @Created at 2021-02-02
// @Modified at 2021-02-02
package router

import (
	v1 "candidate/controller/api/v1"
	_ "candidate/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 候补服务
// @version 1.0
// @description 负责处理与支付和退款相关的业务

// @contact.name LiuYong
// @contact.email ly_yong@qq.com

// @host localhost:8102
// @BasePath /candidate/api/v1
// @query.collection.format multi
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiV1 := r.Group("/candidate/api/v1")
	{
		apiV1.POST("/", v1.Candidate)
		apiV1.POST("/cash", v1.Cash)
		apiV1.POST("/state", v1.ReadState)
	}
	return r
}
