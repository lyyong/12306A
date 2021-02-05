/*
* @Author: 余添能
* @Date:   2021/2/4 10:31 下午
 */
package v1

import (
	"12306A/server/search/rdb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func QueryAllStation(c *gin.Context)  {
	stations:=rdb.QueryStation()

	fmt.Println(stations)
	c.JSON(http.StatusOK,gin.H{"stations":stations})
	//c.String(http.StatusOK,"stations",stations)
	//c.HTML(http.StatusOK,"stations",stations)
}
