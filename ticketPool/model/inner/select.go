/*
* @Author: 余添能
* @Date:   2021/1/24 8:01 下午
 */
package inner

//查找产生的中间类型

type SuitTrainNo struct {
	InitialTime string //无法将datetime按照time.Time类型读出
	TrainNo     string
	StartTime   string //上车时间
	EndTime     string //下车时间
}

//车厢
type Carriage struct {
	ID            int
	TrainNo       string
	CarriageNo    int
	CarriageClass string
	SeatNum       int
}

//座位
type SeatPool struct {
	ID           int
	TicketPoolID int
	TrainNo      string
	SeatNo       string
}

type Seat struct {
}
