/*
* @Author: 余添能
* @Date:   2021/2/23 2:09 下午
 */
package dao

import (
	"fmt"
	"ticketPool/model/inner"
)

func QueryCarriageTypesAll() []*inner.CarriageType {
	sqlStr:="select soft_berth_number,hard_berth_number,senior_soft_benth_number," +
		"hard_seat_number,second_seat_number,first_seat_number,business_seat_number," +
		"business_seat,first_seat,second_seat,hard_seat,hard_berth,soft_berth,senior_soft_berth " +
		"from carriage_types;"
	rows,err:= Db.Query(sqlStr)
	if err!=nil{
		fmt.Println("query table carriage_types failed,err:",err)
		return nil
	}
	var carriageTypes []*inner.CarriageType
	for rows.Next(){
		carriageType:=&inner.CarriageType{}
		rows.Scan(&carriageType.SoftBenthNumber,&carriageType.HardBenthNumber,&carriageType.SeniorSoftBenthNumber,
			&carriageType.HardBenthNumber,&carriageType.SecondSeatNumber,&carriageType.FirstSeatNumber,&carriageType.BusinessSeatNumber,
			&carriageType.BusinessSeat,&carriageType.FirstSeat,&carriageType.SecondSeat,&carriageType.HardSeat,
			&carriageType.HardBenth,&carriageType.SoftBench,&carriageType.SeniorSoftBenth)
		carriageTypes=append(carriageTypes,carriageType)
	}
	return carriageTypes
}
