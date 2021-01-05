package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"user/utils/resp"
)

// Register godoc
// @Summary 用户注册
// @Description 用户注册，参数为用户名和密码
// @ID register-by-username-password
// @Accept mpfd
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} resp.Response
// @Router /register [post]
func Register(c *gin.Context) {
	//TODO: 注册
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println("用户名:", username)
	fmt.Println("密码:", password)

	r := struct {
		Hello    string
		Username string
		Password string
	}{
		Hello:    "Hello，注册成功！",
		Username: username,
		Password: password,
	}
	c.JSON(http.StatusOK, resp.R(r))
}
