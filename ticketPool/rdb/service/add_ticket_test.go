/*
* @Author: 余添能
* @Date:   2021/2/25 9:09 下午
 */
package service

import (
	"testing"
	"ticketPool/model/outer"
)

func TestAddTicket(t *testing.T) {
	ticket:=&outer.Ticket{}
	ticket.TrainNumber="G21"
	ticket.StartTime="2021-02-25 00:00:00"
	ticket.StartStation="北京南"
	ticket.StartStationNum="5"
	ticket.EndTime="2021-02-25 05:14:00"
	ticket.EndStation="上海虹桥"
	ticket.EndStationNum="100"
	ticket.SeatClass="secondSeat"
	ticket.CarriageNum="300"
	ticket.SeatNum="10A"

	AddTicket(ticket)
}
