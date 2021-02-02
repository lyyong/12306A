/*
* @Author: 余添能
* @Date:   2021/1/31 12:26 上午
 */
package init_data

////根据车次去查找具体车厢号，车厢等级，票价等，每趟车次30站，共1900万左右条记录
////目前只写入30车次作为测试
//func WriteTicketPool() {
//	trains := ReadTotalTrainNo()
//
//	sqlStr := "insert into ticket_pool(train_no,initial_time,start_time,start_station,end_time,end_station,carriage_no,carriage_class,ticket_price) " +
//		"values(?,?,?,?,?,?,?,?,?);"
//	sqlStr2 := "insert into ticket_pool(train_no,initial_time,start_time,start_station,end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price) " +
//		"values(?,?,?,?,?,?,?,?,?,?);"
//	st, err := Db.Prepare(sqlStr)
//	st2,err:=Db.Prepare(sqlStr2)
//	defer st.Close()
//	defer st2.Close()
//
//	if err != nil {
//		fmt.Println("prepare table ticket_pool failed, err:", err)
//		return
//	}
//
//	fmt.Println(len(trains))
//	num:=0
//	for _, train := range trains {
//
//		num++
//		if num>30{
//			break
//		}
//		stations := train.Stations
//		n := len(stations)
//		//所有车30节车厢
//		carriageNum := 30
//		for i := 0; i < n; i++ {
//			startStation := stations[i]
//			for j := i + 1; j < n; j++ {
//				endStation := stations[j]
//
//				if i==0 && j==n-1{
//					if strings.Contains(train.TrainNo, "G") {
//						//1-5节商务座/软卧，6-15节一等座，16-30节二等座
//						for k := 1; k <= 5; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 3
//							st2.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "商务座",100, ticketPrice)
//						}
//						for k := 6; k <= 10; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 2
//							st2.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "一等座",100, ticketPrice)
//						}
//						for k := 11; k <= carriageNum; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 1
//							st2.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "二等座",100, ticketPrice)
//						}
//					} else {
//						//1-5节商务座/软卧，6-15节一等座，16-30节二等座
//						for k := 1; k <= 5; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 3
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "软卧", ticketPrice)
//						}
//						for k := 6; k <= 10; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 2
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "硬卧", ticketPrice)
//						}
//						for k := 11; k <= carriageNum; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 1
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "硬座", ticketPrice)
//						}
//					}
//				}else{
//					if strings.Contains(train.TrainNo, "G") {
//						//1-5节商务座/软卧，6-15节一等座，16-30节二等座
//						for k := 1; k <= 5; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 3
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "商务座", ticketPrice)
//						}
//						for k := 6; k <= 10; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 2
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "一等座", ticketPrice)
//						}
//						for k := 11; k <= carriageNum; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 1
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "二等座", ticketPrice)
//						}
//					} else {
//						//1-5节商务座/软卧，6-15节一等座，16-30节二等座
//						for k := 1; k <= 5; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 3
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "软卧", ticketPrice)
//						}
//						for k := 6; k <= 10; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 2
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "硬卧", ticketPrice)
//						}
//						for k := 11; k <= carriageNum; k++ {
//							ticketPrice := (endStation.Price - startStation.Price) * 1
//							st.Exec(train.TrainNo, train.InitialTime, startStation.DepartTime,
//								startStation.StationName, endStation.ArriveTime, endStation.StationName, k, "硬座", ticketPrice)
//						}
//					}
//				}
//
//
//
//			}
//		}
//	}
//}
//
