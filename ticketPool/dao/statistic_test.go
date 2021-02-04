/*
* @Author: 余添能
* @Date:   2021/1/25 9:48 下午
 */
package dao

import (
	"fmt"
	"testing"
)

func TestCountSubTicketPool(t *testing.T) {
	CountSubTicketPool(0)
	num := 0
	for i := 0; i < 100; i++ {
		num += CountSubTicketPool(i)
	}
	fmt.Println("总数：", num)
}
