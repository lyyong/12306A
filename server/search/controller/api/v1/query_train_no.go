/*
* @Author: 余添能
* @Date:   2021/2/4 10:56 下午
 */
package v1

import (
	"12306A/server/search/model/outer"
	"12306A/server/search/rdb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
//http://localhost:18081/search/api/v1/remainder?date=2021-1-23%202001:02:03&startCity=%E5%8C%97%E4%BA%AC&endCity=%E4%B8%8A%E6%B5%B7
func Query(c *gin.Context)  {
	search:=&outer.Search{}
	search.Date=c.Query("date")
	search.StartCity=c.Query("startCity")
	search.EndCity=c.Query("endCity")
	fmt.Println(search)
	trains:=rdb.Query(search)
	if trains==nil{
		c.JSON(http.StatusNoContent,gin.H{
			"code":http.StatusNoContent,
			"msg":"车次票数",
			"data":gin.H{"list":""},
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code":http.StatusOK,
		"msg":"车次票数",
		"data":gin.H{"list":trains},
	})
}
