/*
* @Author: 余添能
* @Date:   2021/2/26 3:16 上午
 */
package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"ticketPool/tools/setting"
)

var(
	Db *sql.DB
	err error
	TrainIds map[string]uint32
	TrainNumbers map[uint32]string
	StationIds map[string]uint32
	StationNames map[uint32]string
	SeatTypes map[uint32]string
	SeatTypeIds map[string]uint32

)

func InitDB()  {
	//Db, err = sql.Open("mysql", "root:12345678@tcp(localhost:3306)/12306a_test")
	////fmt.Println(Db)
	Db, err = sql.Open("mysql", setting.DB.Username+":"+setting.DB.Password+"@"+
		"tcp("+setting.DB.Host+")/"+setting.DB.DbName)
	//Db,err:=sql.Open("mysql",setting.D)
	Db.SetMaxOpenConns(0)
	//fmt.Println(Db)
	if err != nil {
		panic(err.Error())
	}
	InitId()
}
func InitId()  {
	TrainIds =make(map[string]uint32)
	TrainNumbers=make(map[uint32]string)
	StationNames=make(map[uint32]string)
	StationIds = make(map[string]uint32)
	SeatTypes=make(map[uint32]string)
	SeatTypeIds=make(map[string]uint32)
	stations:= SelectStationAll()
	for _,s:=range stations{
		StationIds[s.Name]=uint32(s.ID)
		StationNames[uint32(s.ID)]=s.Name
	}
	//for _,v:=range StationIds {
	//	fmt.Println(v)
	//}
	trains:= SelectTrainAll()
	for _,t:=range trains{
		TrainIds[t.Number]=uint32(t.ID)
		TrainNumbers[uint32(t.ID)]=t.Number
	}

	SeatTypes[0]="businessSeat"
	SeatTypes[1]="firstSeat"
	SeatTypes[2]="secondSeat"
	SeatTypeIds["businessSeat"]=0
	SeatTypeIds["firstSeat"]=1
	SeatTypeIds["secondSeat"]=2
	//for _,v:=range TrainIds {
	//	fmt.Println(v)
	//}
}

func GetTrainNumber(trainId uint32) string  {
	return TrainNumbers[trainId]
}
func GetTrainId(trainNumber string) uint32  {
	return TrainIds[trainNumber]
}

func GetStationName(stationId uint32) string  {
	return StationNames[stationId]
}
func GetStationId(stationName string) uint32  {
	return StationIds[stationName]
}


func GetSeatType(seatId uint32) string  {
	return SeatTypes[seatId]
}

func GetSeatTypeId(seatType string) uint32 {
	return SeatTypeIds[seatType]
}
