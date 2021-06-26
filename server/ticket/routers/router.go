package routers

import (
	"common/middleware/token/usertoken"
	"ticket/controller"
	"ticket/utils/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(setting.Server.RunMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// 使用鉴权中间件
	r.Use(usertoken.TokenParser())
	v1 := r.Group("/ticket/api/v1")
	{
		v1.POST("/buy", controller.BuyTicket)
		v1.POST("/refund", controller.RefundTicket)
		v1.POST("/change", controller.ChangeTicket)
	}

	return r
}
