/**
 * @Author fzh
 * @Date 2020/2/1
 */
package router

import (
	"common/middleware/token/usertoken"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "user/docs"
	"user/router/user"
	"user/util/validator"
)

// 初始化路由
// 如果项目包含很多模块，在这个函数中分别初始化
func InitRouter() *gin.Engine {
	r := gin.New()

	// Token解析中间件
	r.Use(usertoken.TokenParser())

	// 初始化验证器
	validator.InitValidator()

	// 初始化各个模块的路由
	InitUserRouter(r)

	return r
}

// 用户管理模块路由
// @title 用户管理 API
// @BasePath /user/api/v1
func InitUserRouter(r *gin.Engine) *gin.RouterGroup {
	g := r.Group("/user/api/v1")

	// 在路由中添加Swagger
	SwaggerRouter(g)

	// 用户路由配置
	user.Router(g)

	return g
}

// 配置每个模块的Swagger路由
func SwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
