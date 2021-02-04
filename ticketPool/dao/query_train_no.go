/*
* @Author: 余添能
* @Date:   2021/1/24 2:41 下午
 */
package dao

import (
	"12306A/ticketPool/init_data"
	"12306A/ticketPool/model/inner"
	"fmt"
	"time"
)

//根据日期和城市 从train_no中找到所有满足的车次
func SelectTrainNoByCity(startCity, endCity string, startTime time.Time) []*inner.SuitTrainNo {
	//找出还未开的车次，已经过了时间的忽略
	sqlStr := "select initial_time, train_no from train_pool " +
		"where start_time>=? and start_city=? and end_city=?;"
	st, err := init_data.Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("select train_pool failed, err:", err)
		return nil
	}
	rows, err := st.Query("2021-01-23 06:00:00", startCity, endCity)
	if err != nil {
		fmt.Println(err)
	}
	var suitTrainNos []*inner.SuitTrainNo
	for rows.Next() {
		suitTrainNo := &inner.SuitTrainNo{}
		err := rows.Scan(&suitTrainNo.InitialTime, &suitTrainNo.TrainNo)
		if err != nil {
			fmt.Println("scan SuitTrainNo failed, err:", err)
			return nil
		}
		//fmt.Println(initialTime,trainNo)
		suitTrainNos = append(suitTrainNos, suitTrainNo)
	}
	rows.Close()
	return suitTrainNos
}

//根据车次查找各车厢等级的座位数
func SelectTicketNumByTrainNo(trainNo *inner.SuitTrainNo) {

	//统计商务座、一等座、二等座的数量
	//sqlStr:="select  carriage_class,count(seat_no) " +//找到所有座位数
	//	"from seat_pool as seat, " +
	//	"(select id,carriage_class from ticket_pool " + //查询到该车次所有车厢
	//	"where initial_time=? and train_no=? and ) as carriage " +
	//	"where seat.ticket_pool_id=carriage.id " +
	//	"group by carriage_class;"//对车厢类型分组
	//start_time>=?过滤掉已经出发的车次
	//end_time>=?更长的票也算合适的票
	sqlStr := "select  carriage_class,count(seat_no) " +
		"from seat_pool as seat, " +
		"(select carriage_class,count(seat_no) " +
		"from ticket_pool as t " +
		"initial_time=? and train_no=? and start_time>? and end_time>=? ) " +
		""

	st, err := init_data.Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare select ticketNumByTrainNo failed, err:", err)
		return
	}
	rows, _ := st.Query(trainNo.InitialTime, trainNo.TrainNo)
	var seatNums []*inner.Carriage
	for rows.Next() {
		seatNum := &inner.Carriage{}
		rows.Scan(&seatNum.TrainNo, &seatNum.CarriageClass, seatNum.SeatNum)
		seatNums = append(seatNums, seatNum)
		fmt.Println(seatNum)
	}

}
func SelectCarriageByCityAndTrainNo() {

	//从train_pool查到所有满足条件的车次，再去
	sqlStr := "select id,train_no,carriage_class,carriage_no " +
		"from (select initial_time, train_no " +
		"from train_pool " +
		"where start_time>startTime and start_city=startCity and end_city=endCity) as train," +
		"ticket_pool as ticket " +
		"where train.initial_time=ticket.initialTime and train_no=trainNo; "
	rows, err := init_data.Db.Query(sqlStr)
	if err != nil {
		fmt.Println("select ticket failed, err:", err)
		return
	}
	var carriages []*inner.Carriage
	for rows.Next() {
		carriage := &inner.Carriage{}
		rows.Scan(&carriage.ID, carriage.TrainNo, carriage.CarriageNo, carriage.CarriageClass)
		carriages = append(carriages, carriage)
	}

	for _, v := range carriages {
		fmt.Println(v)
	}

}
func InsertTicket() {

}
