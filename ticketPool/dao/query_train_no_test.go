/*
* @Author: 余添能
* @Date:   2021/1/24 5:02 下午
 */
package dao

import (
	"12306A/ticketPool/model/inner"
	"fmt"
	"testing"
	"time"
)

func TestSelectTrainNoByCity(t *testing.T) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-01-23 00:00:00", time.Local)
	suitTrainNos := SelectTrainNoByCity("上海", "北京", startTime)
	for _, suitTrainNo := range suitTrainNos {
		fmt.Println(suitTrainNo)
	}
	fmt.Println(len(suitTrainNos))

}

func TestSelectTicketNumByTrainNo(t *testing.T) {
	suitTrainNo := &inner.SuitTrainNo{
		InitialTime: "2021-01-23 00:00:00",
		TrainNo:     "K4606",
	}
	SelectTicketNumByTrainNo(suitTrainNo)
}
