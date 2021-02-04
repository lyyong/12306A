package controller

import (
	"common/tools/logging"
	"github.com/gin-gonic/gin"
	"indent/service"
	"net/http"
	"strconv"
)

func GetAllIndent (c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {

	}
	indents, err := service.GetAllIndent(userId)
	if err != nil {
		logging.Error("get indent err:", err)
		return
	}
	// 返回 indents 的 json格式
	c.JSON(http.StatusOK, gin.H{"data": indents})
}
