/*
* @Author: 余添能
* @Date:   2021/2/4 10:57 下午
 */
package main

import (
	"12306A/server/search/router"
	"fmt"
)

func main()  {
	r:=router.InitRouter()
	fmt.Println(r)
}
