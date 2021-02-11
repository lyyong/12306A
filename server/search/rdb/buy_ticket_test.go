/*
* @Author: 余添能
* @Date:   2021/2/4 8:14 下午
 */
package rdb

import (
	"12306A/server/search/model/outer"
	"testing"
)

func TestBuyTicketByTrainNoAndDate(t *testing.T) {
	buyTicket:=&outer.BuyTicket{}
	buyTicket.Date="2021-1-23 00:00:00"
	buyTicket.TrainNo="Z4515"
	buyTicket.StartStation="哈尔滨"
	buyTicket.EndStation="南京"
	buyTicket.SeatClass="secondSeat"
	BuyTicketByTrainNoAndDate(buyTicket)
}
