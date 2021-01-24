// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package controller

import (
	"github.com/gin-gonic/gin"
	"pay/tools/message"
)

type Send struct {
	c *gin.Context
}

func (s *Send) Response(httpCode int, errorCode int, data interface{}) {
	s.c.JSON(
		httpCode,
		gin.H{
			"code": errorCode,
			"msg":  message.GetMsg(errorCode),
			"data": data,
		},
	)
}

func NewSend(c *gin.Context) *Send {
	return &Send{c: c}
}
