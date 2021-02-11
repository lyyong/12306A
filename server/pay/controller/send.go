// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-02-25
package controller

import (
	"github.com/gin-gonic/gin"
	"pay/tools/message"
)

type Send struct {
	c *gin.Context
}

type JSONResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (s *Send) Response(httpCode int, result *JSONResult) {
	s.c.JSON(
		httpCode,
		result,
	)
}

func NewSend(c *gin.Context) *Send {
	return &Send{c: c}
}

func NewJSONResult(errorCode int, data interface{}) *JSONResult {
	return &JSONResult{
		Code: errorCode,
		Msg:  message.GetMsg(errorCode),
		Data: data,
	}
}
