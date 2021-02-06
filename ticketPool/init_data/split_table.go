/*
* @Author: 余添能
* @Date:   2021/1/25 6:03 下午
 */
package init_data

import (

	_ "strconv"
)

//
//func SplitWriteTicketPool(ch chan []string, wait sync.WaitGroup) {
//
//	var wg sync.WaitGroup
//	for i := 0; i < 100; i++ {
//		wg.Add(1)
//		go func() {
//			for {
//				columns, s := <-ch
//				if s == false {
//					break
//				}
//				hashValue := crc32.ChecksumIEEE([]byte(columns[1]))
//				tableId := hashValue % (uint32)(TICKETPOOLNUM)
//				a := int(tableId)
//				qry := "insert into ticket_pool_" + strconv.Itoa(a) + " (id,train_no,initial_time,start_time,start_station," +
//					"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price) " +
//					"values (?,?,?,?,?,?,?,?,?,?,?);"
//				//fmt.Println(qry)
//				_, err = Db.Exec(qry, columns[0], columns[1], columns[2], columns[3], columns[4],
//					columns[5], columns[6], columns[7], columns[8], columns[9], columns[10])
//				//if err != nil {
//				//	fmt.Println("exec failed ",err)
//				//}
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//	wait.Done()
//}
//
////func MakeHashTrans()  {
////	sts=make([]*sql.Stmt,100)
////
////	for i:=0;i<100;i++{
////		sqlStr:="insert into ticket_pool_"+strconv.Itoa(i)+" (id,train_no,initial_time,start_time,start_station," +
////			"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price) " +
////			"values(?,?,?,?,?,?,?,?,?,?,?);"
////		sts[i],err=Db.Prepare(sqlStr)
////		if err!=nil{
////			fmt.Println(err)
////			return
////		}
////	}
////	rows,err := Db.Query("select id,train_no,initial_time,start_time,start_station," +
////		"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price  from ticket_pool ;")
////	if err != nil {
////		fmt.Println(err)
////	}
////
////
////	var wait sync.WaitGroup
////	ch:=make(chan []string,50000)
////	wait.Add(1)
////	go SplitWriteTicketPool(ch,wait)
////	for rows.Next() {
////		columns,_:=rows.Columns()
////		err  := rows.Scan(&columns[0], &columns[1], &columns[2], &columns[3],
////			&columns[4],&columns[5],&columns[6],&columns[7],
////			&columns[8],&columns[9],&columns[10])
////		if err!=nil{
////			fmt.Println("scan failed")
////		}
////		ch<-columns
////		//fmt.Println(columns)
////	}
////	close(ch)
////	wait.Wait()
////
////}
//
//func WriteToSubTableFromTicketPool() {
//
//	rows, err := Db.Query("select id,train_no,initial_time,start_time,start_station," +
//		"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price from ticket_pool;")
//	if err != nil {
//		fmt.Println("query ticket_pool failed, err", err)
//		return
//	}
//
//	var wait sync.WaitGroup
//	ch := make(chan []string, 50000)
//	wait.Add(1)
//	go SplitWriteTicketPool(ch, wait)
//	for rows.Next() {
//		columns, _ := rows.Columns()
//		err := rows.Scan(&columns[0], &columns[1], &columns[2], &columns[3],
//			&columns[4], &columns[5], &columns[6], &columns[7],
//			&columns[8], &columns[9], &columns[10])
//		if err != nil {
//			fmt.Println("scan failed")
//		}
//		ch <- columns
//	}
//	close(ch)
//	wait.Wait()
//	//for rows.Next(){
//	//	columns,_:=rows.Columns()
//	//	fmt.Println(columns)
//	//	err := rows.Scan(&columns[0], &columns[1], &columns[2], &columns[3],
//	//		&columns[4], &columns[5], &columns[6], &columns[7],
//	//		&columns[8], &columns[9], &columns[10])
//	//
//	//	if err != nil {
//	//		fmt.Println("scan failed ", err)
//	//	}
//	//
//	//	hashValue := crc32.ChecksumIEEE([]byte(columns[1]))
//	//	tableId := hashValue % (uint32)(TICKETPOOLNUM)
//	//	a := int(tableId)
//	//
//	//	qry := "insert into ticket_pool_"+strconv.Itoa(a)+" (id,train_no,initial_time,start_time,start_station," +
//	//		"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price) " +
//	//		"values (?,?,?,?,?,?,?,?,?,?,?);"
//	//	fmt.Println(qry)
//	//	_, err = Db.Exec(qry,columns[0],columns[1],columns[2],columns[3],columns[4],
//	//		columns[5],columns[6],columns[7],columns[8],columns[9],columns[10])
//	//	if err != nil {
//	//		fmt.Println("exec failed ",err)
//	//	}
//	//}
//	//
//	//left:=1
//	//right:=50000
//	//for {
//	//
//	//
//	//	rows,err := Db.Query("select id,train_no,initial_time,start_time,start_station," +
//	//		"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price  from ticket_pool where " +
//	//		"id>="+strconv.Itoa(left)+" and id<="+strconv.Itoa(right)+";")
//	//	if err != nil {
//	//		fmt.Println(err)
//	//		break
//	//	}
//	//
//	//	for rows.Next() {
//	//		columns, _ := rows.Columns()
//	//		err := rows.Scan(&columns[0], &columns[1], &columns[2], &columns[3],
//	//			&columns[4], &columns[5], &columns[6], &columns[7],
//	//			&columns[8], &columns[9], &columns[10])
//	//		if err != nil {
//	//			fmt.Println("scan failed ", err)
//	//		}
//	//
//	//		hashValue := crc32.ChecksumIEEE([]byte(columns[1]))
//	//		tableId := hashValue % (uint32)(TICKETPOOLNUM)
//	//		a := int(tableId)
//	//		sqlStr:="insert into ticket_pool_"+strconv.Itoa(a)+" (id,train_no,initial_time,start_time,start_station," +
//	//			"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price) " +
//	//			"values(?,?,?,?,?,?,?,?,?,?,?);"
//	//		st,err:=Db.Prepare(sqlStr)
//	//		_,err=st.Exec(columns[0],columns[1],columns[2],columns[3],columns[4],columns[5],
//	//			columns[6],columns[7],columns[8],columns[9],columns[10])
//	//		//qry := "insert into ticket_pool_"+strconv.Itoa(a)+" (id,train_no,initial_time,start_time,start_station," +
//	//		//	"end_time,end_station,carriage_no,carriage_class,seat_num,ticket_price) " +
//	//		//	"values ("+columns[0]+","+columns[1]+","+columns[2]+","+ columns[3]+","+
//	//		//	columns[4]+","+columns[5]+","+columns[6]+","+columns[7]+","+columns[8]+","+columns[9]+","+columns[10]+");"
//	//		//fmt.Println(qry)
//	//		//_, err = Db.Exec(qry)
//	//		//if err != nil {
//	//		//	fmt.Println("exec failed ",err)
//	//		//}
//	//
//	//	}
//	//	left+=50000
//	//	right+=50000
//
//	//}
//
//}
