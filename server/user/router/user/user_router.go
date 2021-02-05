/**
 * @Author fzh
 * @Date 2020/2/1
 */
package user

import (
	"github.com/gin-gonic/gin"
	"user/api/user"
)

func Router(r *gin.RouterGroup) *gin.RouterGroup {
	// 测试接口
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "测试成功")
	})

	r.POST("/register", user.Register)
	r.POST("/login", user.Login)
	return r
}
