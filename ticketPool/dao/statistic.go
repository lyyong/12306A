/*
* @Author: 余添能
* @Date:   2021/1/25 9:41 下午
 */
package dao

import (
	"12306A/ticketPool/init_data"
	"fmt"
	"strconv"
)

func CountSubTicketPool(i int) int {
	sqlStr := "select count(*) from ticket_pool_" + strconv.Itoa(i) + ";"

	row, err := init_data.Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var count int
	for row.Next() {
		row.Scan(&count)
		fmt.Println(count)
	}
	return count
}
