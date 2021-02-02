/*
* @Author: 余添能
* @Date:   2021/1/23 10:11 下午
 */
package inner

import "time"

//通过坐车时间和城市去找合适的车次
type TrainPool struct {
	ID int
	//InitialTime  time.Time
	//TerminalTime time.Time
	TrainNo   string
	StartCity string
	//StartStation string
	StartTime string
	EndCity   string
	//EndStation   string
	EndTime string
}

type TicketPool struct {
	ID            int
	TrainNo       string
	InitialTime   string
	StartTime     time.Time //上车时间
	StartStation  string
	EndTime       time.Time
	EndStation    string
	CarriageNo    int     //车厢号
	CarriageClass string  //车厢等级
	SeatNo        int     //车厢内可用座位数
	TicketPrice   float64 //票价
}

//
////从车次去找票
//type TicketPool struct {
//	ID            int
//	TrainNo       string
//	InitialTime   time.Time //起始时间
//	StartTime     time.Time //上车时间
//	StartStation  string
//	EndTime       time.Time
//	EndStation    string
//	CarriageNo    int     //车厢号
//	CarriageClass string  //车厢等级
//	SeatNum       int     //车厢内可用座位数
//	TicketPrice   float64 //票价
//}
