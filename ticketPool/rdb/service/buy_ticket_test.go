/*
* @Author: 余添能
* @Date:   2021/2/25 8:50 下午
 */
package service

import (
	"fmt"
	"testing"
	"ticketPool/model/outer"
)

func TestBuyTicket(t *testing.T) {
	var buyTickets []*outer.BuyTicket
	buyTicket:=&outer.BuyTicket{}
	buyTicket.TrainNumber="G21"
	buyTicket.StartStation="天津南"
	buyTicket.StartTime="2021-02-25 00:00:00"
	buyTicket.EndStation="济南西"
	buyTicket.EndTime="2021-02-25 05:14:00"
	buyTicket.SeatClass="secondSeat"
	buyTicket.SeatPlace=""
	buyTickets=append(buyTickets,buyTicket)
	tickets:=BuyTicket(buyTickets)
	for _,v:=range tickets{
		fmt.Println(v)
	}
}
