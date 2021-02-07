/*
* @Author: 余添能
* @Date:   2021/1/31 12:03 上午
 */
package init_data

import (
	"12306A/ticketPool/model/inner"
	"fmt"
)

//初始化train_pool，用于查找两个城市之间的车次
//共有36万条记录
func WriteTrainPool() {
	fmt.Println("开始初始化train_pool表")
	trains := ReadTotalTrainNo()

	sqlStr := "insert into train_pool(initial_time,terminal_time,train_no,start_city,start_time,end_city,end_time) " +
		"values(?,?,?,?,?,?,?);"
	st, err := Db.Prepare(sqlStr)
	defer st.Close()
	if err != nil {
		fmt.Println("prepare table train_pool failed, err:", err)
		return
	}

	for _, train := range trains {

		n := len(train.Stations)
		stations := train.Stations
		//记录上车城市是否写入
		startCityWrited := make(map[string]string)
		for i := 0; i < n; i++ {
			//上车城市
			startCity := stations[i].CityName
			if startCityWrited[startCity] != "" {
				continue
			}
			startCityWrited[startCity] = startCity
			//记录下车城市是否写入
			endCityWrited := make(map[string]string)
			for j := i + 1; j < n; j++ {
				//下车城市
				endCity := stations[j].CityName
				//去重
				//一趟车可能会经过同一个城市的两个车站，比如下属县级市，不重复输入
				if endCityWrited[endCity] == "" {
					//车次的：起始时间，车次终止时间，车次，上车城市，上车时间，下车城市，下车时间
					st.Exec(train.InitialTime, train.TerminalTime, train.TrainNo, startCity, stations[i].ArriveTime, endCity, stations[j].ArriveTime)
					endCityWrited[endCity] = endCity
				}
			}
		}
	}
}

func ReadTrainPoolAll() []*inner.TrainPool {
	strSql := "select train_no,start_city,start_time,end_city,end_time from train_pool;"
	rows, err := Db.Query(strSql)
	if err != nil {
		fmt.Println("select train_pool failed, err:", err)
		return nil
	}
	var trainPools []*inner.TrainPool
	for rows.Next() {
		trainPool := &inner.TrainPool{}
		rows.Scan(&trainPool.TrainNo, &trainPool.StartCity, &trainPool.StartTime, &trainPool.EndCity, &trainPool.EndTime)
		trainPools = append(trainPools, trainPool)
	}
	fmt.Println(len(trainPools))
	return trainPools
}
