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

func Query(c *gin.Context)  {
	search:=&outer.Search{}
	search.Date=c.Param("date")
	search.StartCity=c.Param("startCity")
	search.EndCity=c.Param("endCity")
	trains:=rdb.Query(search)
	if trains==nil{
		return
	}
	fmt.Println(trains)
	c.JSON(http.StatusOK,gin.H{"trains":trains})
}
