package routers

import (
	"common/middleware/token/usertoken"
	"github.com/gin-gonic/gin"
	"ticket/controller"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 使用鉴权中间件
	r.Use(usertoken.TokenParser())
	v1 := r.Group("/ticket/api/v1")
	{
		v1.POST("/buyTicket", controller.BuyTicket)
		v1.POST("/RefundTicket", controller.RefundTicket)
	}

	return r
}
