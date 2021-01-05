package controller

import (
	"12306A/server/pay/tools/message"
	"github.com/gin-gonic/gin"
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
