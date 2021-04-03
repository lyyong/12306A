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
	r.POST("/buyTicket", controller.BuyTicket)
	r.POST("/RefundTicket", controller.RefundTicket)
	return r
}
