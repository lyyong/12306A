/*
* @Author: 余添能
* @Date:   2021/2/24 1:58 下午
 */
package dao

import (
	"fmt"
	"ticketPool/model/outer"
	"time"
)

func InsertTicket(ticket *outer.Ticket) bool {
	sqlStr:="insert into ticket_pools(created_at,created_by,date,train_number,start_time,start_station,start_station_num,end_time,end_station,end_station_num," +
		"seat_class,carriage_num,seat_num,price) values(?,?,?,?,?,?,?,?,?,?,?,?,?);"
	_,err:=Db.Exec(sqlStr,time.Now(),"余添能",ticket.Date,ticket.TrainNumber,ticket.StartTime,ticket.StartStation,ticket.StartStationNum,
		ticket.EndTime,ticket.EndStation,ticket.EndStationNum,ticket.SeatClass,ticket.CarriageNum,ticket.SeatNum,ticket.Price)
	if err!=nil{
		fmt.Println("insert ticket_pools failed, err:",err)
		return false
	}
	return true
}
