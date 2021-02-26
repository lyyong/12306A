/*
* @Author: 余添能
* @Date:   2021/2/25 6:55 下午
 */
package dao

import (
	"fmt"
	"testing"
)

func TestSelectStationAll(t *testing.T) {
	stations:=SelectStationAll()

	for _,s:=range stations{
		fmt.Println(s)
	}
}
