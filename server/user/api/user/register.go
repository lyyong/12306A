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

type RegisterJSON struct {
	Username string
	Password string
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册，参数为用户名和密码
// @ID register-by-username-password
// @Accept json
// @Produce json
// @Param form body RegisterJSON true "注册信息"
// @Success 200 {object} resp.Response
// @Router /register [post]
func Register(c *gin.Context) {
	j := new(RegisterJSON)
	if err := c.ShouldBindJSON(j); err != nil {
		c.JSON(http.StatusBadRequest, resp.R(struct{}{}).SetMsg("JSON格式错误"))
	}

	if err := user.Register(j.Username, j.Password); err != nil {
		if errors.Is(err, errortype.ErrUserNameHasExist) {
			c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户已注册"))
		} else {
			c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("注册失败"))
		}
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("注册成功"))
}
