package user

import (
	"github.com/gin-gonic/gin"
	"user/apis/user"
)

func Router(r *gin.RouterGroup) *gin.RouterGroup {
	// 测试接口
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "测试成功")
	})

	r.POST("/register", user.Register)
	return r
}
