/**
 * @Author fzh
 * @Date 2021/2/6
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

type LoginJSON struct {
	Username string
	Password string
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录，参数为用户名和密码
// @ID login-by-username-password
// @Accept json
// @Produce json
// @Param form body LoginJSON true "登录信息"
// @Success 200 {object} resp.Response
// @Router /login [post]
func Login(c *gin.Context) {
	j := new(LoginJSON)
	if err := c.ShouldBindJSON(j); err != nil {
		c.JSON(http.StatusBadRequest, resp.R(struct{}{}).SetMsg("JSON格式错误"))
	}

	r := struct {
		Token string `json:"token"`
	}{}

	if token, err := user.Login(j.Username, j.Password); err != nil {
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
