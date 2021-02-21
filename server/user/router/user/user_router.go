/**
 * @Author fzh
 * @Date 2020/2/1
 */
package user

import (
	"github.com/gin-gonic/gin"
	"user/api/httpapi"
)

func Router(r *gin.RouterGroup) *gin.RouterGroup {
	// 测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	r.POST("/register", httpapi.Register)
	r.POST("/login", httpapi.Login)
	return r
}
