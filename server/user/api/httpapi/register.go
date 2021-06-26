/**
 * @Author fzh
 * @Date 2020/2/1
 */
package httpapi

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"user/global/errortype"
	"user/service"
	"user/util/resp"
)

type RegisterRequest struct {
	Username          string
	Password          string
	CertificateType   int
	Name              string
	CertificateNumber string `binding:"certificateNumber"`
	PhoneNumber       string `binding:"phoneNumber"`
	Email             string
	PassengerType     int
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册，参数为用户名和密码
// @ID register-by-username-password
// @Accept json
// @Produce json
// @Param form body RegisterRequest true "注册信息"
// @Success 200 {object} resp.Response
// @Router /register [post]
func Register(c *gin.Context) {
	req := new(RegisterRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, resp.R(struct{}{}).SetMsg("JSON格式错误"))
		return
	}

	p := &service.RegisterParam{
		Username:        req.Username,
		Password:        req.Password,
		Name:            req.Name,
		CertificateType: req.CertificateType,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		PassengerType:   req.PassengerType,
	}
	if err := service.Register(p); err != nil {
		if errors.Is(err, errortype.ErrUserNameHasExist) {
			c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("用户已注册"))
		} else {
			c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("注册失败"))
		}
		return
	}

	c.JSON(http.StatusOK, resp.R(struct{}{}).SetMsg("注册成功").SetCode(200))
}
