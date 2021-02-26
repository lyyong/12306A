/*
* @Author: 余添能
* @Date:   2021/2/26 3:16 上午
 */
package dao

var(
	TrainIds map[uint32]string
	StationIds map[uint32]string
	SeatTypes map[uint32]string
	SeatTypeIds map[string]uint32
)

func InitId()  {
	TrainIds =make(map[uint32]string)
	StationIds =make(map[uint32]string)
	SeatTypes=make(map[uint32]string)
	SeatTypeIds=make(map[string]uint32)
	stations:= SelectStationAll()
	for _,s:=range stations{
		StationIds[uint32(s.ID)]=s.Name
	}
	//for _,v:=range StationIds {
	//	fmt.Println(v)
	//}
	trains:= SelectTrainAll()
	for _,t:=range trains{
		TrainIds[uint32(t.ID)]=t.Number
	}

	SeatTypes[1]="businessSeat"
	SeatTypes[2]="firstSeat"
	SeatTypes[3]="secondSeat"
	SeatTypeIds["businessSeat"]=1
	SeatTypeIds["firstSeat"]=2
	SeatTypeIds["secondSeat"]=3
	//for _,v:=range TrainIds {
	//	fmt.Println(v)
	//}
}

func GetTrainNumber(trainId uint32) string  {
	return TrainIds[trainId]
}

func GetStationName(stationId uint32) string  {
	return StationIds[stationId]
}

func GetSeatType(seatId uint32) string  {
	return SeatTypes[seatId]
}

func GetSeatTypeId(seatType string) uint32 {
	return SeatTypeIds[seatType]
}