/*
* @Author: 余添能
* @Date:   2021/1/23 1:10 上午
 */
package init_data

import (
	"12306A/ticketPool/model/inner"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mozillazg/go-pinyin"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func WriteStationProvinceCity() {
	fmt.Println("开始初始化station_prov_city表")
	//先清空表
	Db.Exec("delete from station_province_city")

	bytes, err := ioutil.ReadFile("./ticketPool/stations_prov_city.json")
	if err != nil {
		fmt.Println("read file failed, err:", err)
	}
	//fmt.Println(string(bytes))

	//使用map将json字符串读出来
	var station []map[string]interface{}
	err = json.Unmarshal(bytes, &station)
	if err != nil {
		fmt.Println("json convert struct failed, err:", err)
	}
	//fmt.Println(len(station),"station")

	sqlStr := "insert into station_province_city(province,city,city_code,station_name,station_telecode,station_spell) " +
		"values(?,?,?,?,?,?);"
	st, err := Db.Prepare(sqlStr)
	defer st.Close()
	if err != nil {
		fmt.Println("prepare sql failed, err:", err)
	}
	//获得：城市-城市编码
	cities := ReadCity()
	//写入表
	for _, v := range station {
		city := v["city"].(string)
		city = strings.Trim(city, "市")

		province := v["province"].(string)
		stationName := v["station_name"].(string)
		stationTelecode := v["station_telecode"].(string)
		//生成站点的拼音
		a := pinyin.NewArgs()
		spells := pinyin.Pinyin(stationName, a)
		var stationSpell string
		for _, spell := range spells {
			for _, v := range spell {
				stationSpell += v
			}
		}
		//fmt.Println(stationSpell,province,stationTelecode,stationName)

		_, err := st.Exec(province, city, cities[city], stationName, stationTelecode, stationSpell)
		if err != nil {
			fmt.Println(city, err)
		}
	}
}

//获取城市的编码
func ReadCity() map[string]string {

	bytes, err := ioutil.ReadFile("./ticketPool/city.json")
	if err != nil {
		fmt.Println("read city.json failed, err:", err)
	}
	//fmt.Println(string(bytes))

	var cities map[string]string
	cities = make(map[string]string)
	var data []map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("json convert struct failed, err:", err)
	}
	for _, v := range data {
		//fmt.Println(k,v)

		list := v["list"].([]interface{})
		for _, c := range list {
			t := c.(map[string]interface{})
			//fmt.Println(t["code"],t["label"],t["name"])
			name := t["name"].(string)
			label := t["label"].(string)
			//code:=t["code"].(string)
			cities[name] = label
		}
	}
	fmt.Println("city_num:", len(cities))
	return cities
}

//从station_province_city中获得车站-城市的映射关系
func ReadStationCity() map[string]string {
	sqlStr := "select station_name, city from station_province_city;"
	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println("exec query table station_province_city failed, err:", err)
		return nil
	}
	//获得车站-城市映射关系
	stationCity := make(map[string]string)
	for rows.Next() {
		var station, city string
		rows.Scan(&station, &city)
		stationCity[station] = city
	}
	return stationCity
}

//==================================

//读取车次信息，包括途径所有站点
func ReadTrainNo() []*inner.Train {

	//xlsx, err := excelize.OpenFile("/Users/yutianneng/go/src/12306A/ticketPool/train_no.xlsx")
	xlsx, err := excelize.OpenFile("./ticketPool/train_no.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rows := xlsx.GetRows("Sheet1")

	var trains []*inner.Train
	var num int64 = 0
	for i := 1; ; {
		var stations []*inner.Station
		initial, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-01-23 00:00:00", time.Local)
		train := &inner.Train{}
		train.TrainNo = rows[i][6]
		train.ID = num
		num++
		price := 0.0
		for ; i < len(rows); i++ {
			row := rows[i]
			//另一趟列车
			if row[6] != train.TrainNo {
				break
			}
			sta := &inner.Station{}
			sta.ID = 1
			//票价,一站10元
			sta.Price = price
			price += 10
			//站序
			sta.StationNo, _ = strconv.ParseInt(row[0], 10, 64)
			//站名
			sta.StationName = row[1]
			//里程
			sta.Mileage, _ = strconv.ParseInt(row[5], 10, 64)
			//历时
			//00:00:00没法写入mysql，最小时间
			dura, _ := time.ParseInLocation("2006-01-02 15:04:05", "2001-01-01 01:00:00", time.Local)
			s := strings.Split(row[4], ":")
			hours, _ := time.ParseDuration(s[0] + "h")
			mins, _ := time.ParseDuration(s[1] + "m")
			dura = dura.Add(hours)
			dura = dura.Add(mins)
			sta.Duration = dura
			//fmt.Println(sta.Duration)
			//到站时间
			s1 := strings.Split(row[2], ":")
			hours1, _ := time.ParseDuration(s1[0] + "h")
			mins1, _ := time.ParseDuration(s1[1] + "m")
			arrtime := initial.Add(hours1)
			arrtime = arrtime.Add(mins1)
			sta.ArriveTime = arrtime
			//发车时间
			s2 := strings.Split(row[3], ":")
			hours2, _ := time.ParseDuration(s2[0] + "h")
			mins2, _ := time.ParseDuration(s2[1] + "m")
			starttime := initial.Add(hours2)
			starttime = starttime.Add(mins2)
			sta.DepartTime = starttime
			//不是第一站
			if stations != nil && sta.ArriveTime.Before(stations[len(stations)-1].ArriveTime) {
				//如果当前时间比前一站还小，代表火车横跨了两天
				//Time比较
				d, _ := time.ParseDuration("24h")
				sta.ArriveTime = sta.ArriveTime.Add(d)
				sta.DepartTime = sta.DepartTime.Add(d)
				initial = initial.Add(d)
			}
			stations = append(stations, sta)
		}
		train.Stations = stations
		//总站数
		train.StationNum = int64(len(stations))
		//起始时间
		train.InitialTime = stations[0].ArriveTime
		//停止时间
		train.TerminalTime = stations[len(stations)-1].ArriveTime
		trains = append(trains, train)
		if i == len(rows) {
			break
		}
	}
	return trains
	//一天有12801趟车次
	//for _,v:=range trains{
	//	fmt.Println(v.ID,v.TrainNo,v.StationNum,v.InitialTime,v.TerminalTime)
	//	for _,t:=range v.Stations{
	//		fmt.Println(t)
	//	}
	//}
	//fmt.Println(len(trains))
}

//
//func ReadTicketPoolID(trainNo, startStation, endStation string) []*inner.SeatPool {
//	sqlStr := "select id,train_no from ticket_pool where train_no=? and start_station=? and end_station=?;"
//	st, err := Db.Prepare(sqlStr)
//	defer st.Close()
//
//	if err != nil {
//		fmt.Println("prepare table ticket_pool failed, err:", err)
//		return nil
//	}
//	var ids []*inner.SeatPool
//	rows, _ := st.Query(trainNo, startStation, endStation)
//	for rows.Next() {
//		id := &inner.SeatPool{}
//		rows.Scan(&id.TicketPoolID, &id.TrainNo)
//		ids = append(ids, id)
//	}
//	return ids
//}
//
////每个车厢100个座位，每趟车3000座位
//func WriteSeatPool() {
//
//	////从ticket_pool中查指定车次起点站和终点站
//	//获取全票站点的id,train_no
//	var ticketIds chan *inner.SeatPool
//	ticketIds = make(chan *inner.SeatPool, 400000)
//
//	trains := ReadTrainNo()
//	for _, train := range trains {
//		startStation := train.Stations[0].StationName
//		endStation := train.Stations[len(train.Stations)-1].StationName
//		ids := ReadTicketPoolID(train.TrainNo, startStation, endStation)
//		for _, v := range ids {
//			ticketIds <- v
//		}
//	}
//	fmt.Println(len(ticketIds))
//
//	close(ticketIds)
//	//刚开始，只有全票
//	sqlStr := "insert into seat_pool(train_no,ticket_pool_id,seat_no) values(?,?,?);"
//	st, err := Db.Prepare(sqlStr)
//	defer st.Close()
//	if err != nil {
//		fmt.Println("prepare talbe seat_pool failed, err:", err)
//		return
//	}
//
//	var wg sync.WaitGroup
//	fmt.Println("开始写入seat_pool")
//	for j := 1; j <= 100; j++ {
//		wg.Add(1)
//		go func() {
//			for {
//				t, s := <-ticketIds
//				if s == false {
//					break
//				}
//				for i := 1; i <= 100; i++ {
//					st.Exec(t.TrainNo, t.TicketPoolID, i)
//				}
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}
//
//func WriteSeatNumToTicketPool() {
//	////从ticket_pool中查指定车次起点站和终点站
//	//获取全票站点的id,train_no
//	var ticketIds chan *inner.SeatPool
//	ticketIds = make(chan *inner.SeatPool, 400000)
//
//	trains := ReadTrainNo()
//	for _, train := range trains {
//		startStation := train.Stations[0].StationName
//		endStation := train.Stations[len(train.Stations)-1].StationName
//		ids := ReadTicketPoolID(train.TrainNo, startStation, endStation)
//		for _, v := range ids {
//			ticketIds <- v
//		}
//	}
//	fmt.Println(len(ticketIds))
//
//	close(ticketIds)
//	//刚开始，只有全票
//	sqlStr := "update ticket_pool set seat_num=100 where id=?;"
//	st, err := Db.Prepare(sqlStr)
//	defer st.Close()
//	if err != nil {
//		fmt.Println("prepare table ticket_pool failed, err:", err)
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for j := 1; j <= 100; j++ {
//		wg.Add(1)
//		go func() {
//			for {
//				t, s := <-ticketIds
//				if s == false {
//					break
//				}
//				st.Exec(t.TicketPoolID)
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}
