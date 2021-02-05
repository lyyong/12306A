/**
 * @Author fzh
 * @Date 2020/2/1
 */
package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"user/global/errortype"
	"user/service/user"
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
	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := user.Register(username, password); err != nil {
		if errors.Is(err, errortype.ErrUserNameHasExist) {
			c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户已注册"))
		} else {
			c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("注册失败"))
		}
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("注册成功"))
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录，参数为用户名和密码
// @ID login-by-username-password
// @Accept mpfd
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} resp.Response
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	r := struct {
		Token string `json:"token"`
	}{}

	if token, err := user.Login(username, password); err != nil {
		if errors.Is(err, errortype.ErrUserNotExist) || errors.Is(err, errortype.ErrWrongPassword) {
			c.JSON(http.StatusOK, resp.R(r).SetMsg("用户名或密码不正确"))
		} else {
			c.JSON(http.StatusOK, resp.R(r).SetMsg("登录失败"))
		}
		return
	} else {
		r.Token = token
	}
	c.JSON(http.StatusOK, resp.R(r).SetMsg("登录成功"))
}
