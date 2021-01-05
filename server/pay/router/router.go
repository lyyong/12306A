package router

import (
	v1 "12306A/server/pay/controller/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/payok", v1.PayOK)
	}

	return r
}
