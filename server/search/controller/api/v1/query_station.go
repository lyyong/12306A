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
	if stations==nil{
		c.JSON(http.StatusNoContent,gin.H{
			"code":http.StatusNoContent,
			"stations":"",
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"stations":stations,
		})
	}
	//c.String(http.StatusOK,"stations",stations)
	//c.HTML(http.StatusOK,"stations",stations)
}

func QueryStationByTrainNo(c *gin.Context)  {
	trainNo:=c.Query("train_no")
	stations:=rdb.QueryStationByTrainNo(trainNo)
	if stations==nil{
		c.JSON(http.StatusNoContent,gin.H{
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