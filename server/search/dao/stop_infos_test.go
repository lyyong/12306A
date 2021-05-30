/*
* @Author: 余添能
* @Date:   2021/2/25 10:30 下午
 */
package dao

import (
	"fmt"
	"testing"
)

func TestSelectStopInfoAll(t *testing.T) {
	all := SelectStopInfoAll()
	for _, v := range all {
		fmt.Println(v)
		fmt.Println(v.TrainNumber)
	}
}
