/*
* @Author: 余添能
* @Date:   2021/2/4 11:48 下午
 */
package main

import (
	"fmt"
	"ticketPool/dao"
	"ticketPool/rpc"
)


func main()  {
	dao.InitId()
	rpc.Setup()
	fmt.Println()
}


