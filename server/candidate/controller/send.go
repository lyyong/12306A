// Package controller
// @Author liuYong
// @Created at 2020-01-05
package controller

import (
	"candidate/tools/message"
	"github.com/gin-gonic/gin"
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
