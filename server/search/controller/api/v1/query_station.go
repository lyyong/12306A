/*
* @Author: 余添能
* @Date:   2021/2/4 10:31 下午
 */
package v1

import (
	"12306A-search/dao"
	"12306A-search/rdb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func QueryAllStation(c *gin.Context)  {

	c.JSON(http.StatusOK,gin.H{
		"code":http.StatusOK,
		"msg":"...",
		"data":gin.H{
			"hot_cities":dao.HotCities,
			"cities":dao.CityLists,
		},
	})
	//c.String(http.StatusOK,"stations",stations)
	//c.HTML(http.StatusOK,"stations",stations)
}

func QueryStationByTrainNo(c *gin.Context)  {
	trainNo:=c.Query("train_no")
	fmt.Println("aaa",trainNo)
	stations:=rdb.QueryStationByTrainNo(trainNo)
	if stations==nil{
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusNoContent,
			"data":gin.H{"stations":""},
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"data":gin.H{"stations":stations},
		})
	}
}