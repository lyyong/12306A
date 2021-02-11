/*
* @Author: 余添能
* @Date:   2021/2/3 6:08 下午
 */
package main

import (
	"fmt"
	"strings"
	"time"
)

func main()  {

	now:=time.Now()
	fmt.Println(strings.Split(now.String()," ")[0])

}
