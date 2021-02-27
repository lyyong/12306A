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

func QueryRemainder(c *gin.Context)  {

	date:=c.Query("date")
	startCity:=c.Query("startCity")
	endCity:=c.Query("endCity")
	//fmt.Println(date,startCity,endCity,"aaa")
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
