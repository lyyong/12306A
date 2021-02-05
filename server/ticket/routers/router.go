package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticket/controller"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/buyTicket", controller.BuyTicket)
	r.POST("/testPost", TestPost)
	return r
}

func TestPost(c *gin.Context){
	var btReq controller.BuyTicketRequest
		if err := c.ShouldBindJSON(&btReq); err != nil {
			// 参数有误
			c.JSON(http.StatusOK, controller.Response{
				Code: 0,
				Msg:  "参数有误1",
				Data: nil,
			})
			return
		}
	c.JSON(http.StatusOK, gin.H{"isOk":"OK"})
}