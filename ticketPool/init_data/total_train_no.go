/*
* @Author: 余添能
* @Date:   2021/1/30 10:43 下午
 */
package init_data

import (
	"12306A/ticketPool/model/inner"
	"fmt"
	"time"
)

func WriteTotalTrainNo() {
	sqlStr := "insert into total_train_no(train_no,station_num,initial_time," +
		"terminal_time,station_no,station_name,city_name,arrive_time,depart_time,duration,mileage,price) " +
		"values (?,?,?,?,?,?,?,?,?,?,?,?);"
	st, err := Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare total_train_no failed, err:", err)
		return
	}
	//从文件中读取所有车次数据
	trains := ReadTrainNo()
	//读取站点-城市的映射关系
	stationCity := ReadStationCity()
	for _, train := range trains {

		for _, station := range train.Stations {
			cityName := stationCity[station.StationName]
			_, err := st.Exec(train.TrainNo, train.StationNum, train.InitialTime, train.TerminalTime,
				station.StationNo, station.StationName, cityName, station.ArriveTime, station.DepartTime,
				station.Duration, station.Mileage, station.Price)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

//total_train_no数据有缺失，去除缺失的数据
func ReviseTotalTrainNo() {
	sqlStr := "select train_no from total_train_no where city_name='';"
	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println("query total_train_no failed, err:", err)
		return
	}
	delSqlStr := "delete from total_train_no where train_no=?;"
	st, err := Db.Prepare(delSqlStr)
	if err != nil {
		fmt.Println("prepare total_train_no failed, err:", err)
		return
	}
	//rows中是所有出现缺失数据的车次
	for rows.Next() {
		var trainNo string
		rows.Scan(&trainNo)
		st.Exec(trainNo)
	}
}

func ReadTotalTrainNo() []*inner.Train {
	sqlStr := "select train_no,station_num,initial_time,terminal_time,station_no,station_name," +
		"city_name,arrive_time,depart_time,duration,mileage,price " +
		"from total_train_no; "
	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println("query total_train_no failed, err:", err)
		return nil
	}

	totalTrainNo := make(map[string][]*inner.TotalTrainNo)
	for rows.Next() {
		station := &inner.TotalTrainNo{}
		err := rows.Scan(&station.TrainNo, &station.StationNum, &station.InitialTime, &station.TerminalTime,
			&station.StationNo, &station.StationName, &station.CityName, &station.ArriveTime,
			&station.DepartTime, &station.Duration, &station.Mileage, &station.Price)
		if err != nil {
			fmt.Println(err)
		}
		totalTrainNo[station.TrainNo] = append(totalTrainNo[station.TrainNo], station)
	}
	var trains []*inner.Train
	for _, t := range totalTrainNo {

		train := &inner.Train{}
		stas := t
		train.TrainNo = stas[0].TrainNo
		train.InitialTime, _ = time.ParseInLocation("2006-01-02 15:04:05", stas[0].InitialTime, time.Local)
		train.TerminalTime, _ = time.ParseInLocation("2006-01-02 15:04:05", stas[0].TerminalTime, time.Local)
		train.StationNum = stas[0].StationNum
		var stations []*inner.Station
		for _, station := range stas {
			s := &inner.Station{}
			s.StationNo = station.StationNo
			s.StationName = station.StationName
			s.CityName = station.CityName
			s.ArriveTime, _ = time.ParseInLocation("2006-01-02 15:04:05", station.ArriveTime, time.Local)
			s.DepartTime, _ = time.ParseInLocation("2006-01-02 15:04:05", station.DepartTime, time.Local)
			s.Mileage = station.Mileage
			s.Duration, _ = time.ParseInLocation("2006-01-02 15:04:05", station.Duration, time.Local)
			s.Price = station.Price
			stations = append(stations, s)
		}
		train.Stations = stations
		trains = append(trains, train)
	}
	fmt.Println(len(trains))
	return trains
}
