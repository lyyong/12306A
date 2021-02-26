/*
* @Author: 余添能
* @Date:   2021/2/23 2:30 下午
 */
package dao

import (
	"fmt"
	"testing"
)

func TestQueryCarriageTypesAll(t *testing.T) {
	all := QueryCarriageTypesAll()
	for _,v:=range all{
		fmt.Println(v)
	}
}
