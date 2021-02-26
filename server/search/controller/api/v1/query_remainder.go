/*
* @Author: 余添能
* @Date:   2021/2/4 10:56 下午
 */
package v1

import (
	"12306A-search/rdb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Query(c *gin.Context)  {

	date:=c.PostForm("date")
	startCity:=c.PostForm("startCity")
	endCity:=c.PostForm("endCity")
	//fmt.Println(search)
	//trains:=rdb.Query(search)
	trains:=rdb.QueryTicketNumByDate(date,startCity,endCity)
	if trains==nil{
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusNoContent,
			"msg":"车次票数",
			"data":gin.H{"list":""},
		})
	}else {
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"msg":"车次票数",
			"data":gin.H{"list":trains},
		})
	}

}
