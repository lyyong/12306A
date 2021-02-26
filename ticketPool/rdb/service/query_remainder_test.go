/*
* @Author: 余添能
* @Date:   2021/2/26 12:02 上午
 */
package service

import (
	"fmt"
	"testing"
)

func TestQueryTicketNumByTrainNoAndDate(t *testing.T) {
	date:="2021-02-25"
	trainNo:="G21"
	seatClass:="secondSeat"
	start:="北京南"
	end:="上海虹桥"
	remainderNum:=QueryTicketNumByTrainNoAndDate(date,trainNo,seatClass,start,end)
	fmt.Println(remainderNum)
}
