package routers

import (
	"github.com/gin-gonic/gin"
	"ticket/controller"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/buyTicket", controller.BuyTicket)
	return r
}
