/*
* @Author: 余添能
* @Date:   2021/2/2 3:40 下午
 */
package rdb

import "testing"

func TestQueryTicketNumByTrainNoAndDate(t *testing.T) {
	QueryTicketNumByTrainNoAndDate("2021-1-23","K4729","secondSeat",1,10)
}
