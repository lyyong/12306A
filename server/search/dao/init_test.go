/*
* @Author: 余添能
* @Date:   2021/2/26 12:03 下午
 */
package dao

import (
	"12306A-search/tools/settings"
	"fmt"
	"testing"
)

func TestInitCityLists(t *testing.T) {
	//InitCityLists()
	settings.Setup()
	InitDB()
	fmt.Println(Db)
}

