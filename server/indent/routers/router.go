package routers

import (
	"github.com/gin-gonic/gin"
	"indent/controller"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/getIndent", controller.GetAllIndent)
	r.GET("/getNotStartIndent", controller.GetNotStartIndent)
	return r
}
