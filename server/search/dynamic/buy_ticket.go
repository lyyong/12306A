/*
* @Author: 余添能
* @Date:   2021/2/4 10:29 下午
 */
package dynamic

import (
	"12306A/server/search/model/outer"
	"12306A/server/search/rdb"
)

func BuyTicketByTrainNoAndDate(buyTicket *outer.BuyTicket) *outer.Ticket{
	return rdb.BuyTicketByTrainNoAndDate(buyTicket)
}
