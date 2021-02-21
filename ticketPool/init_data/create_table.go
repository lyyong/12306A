/*
* @Author: 余添能
* @Date:   2021/1/23 1:10 上午
 */
package init_data

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var Db *sql.DB
var err error

func init() {
	Db, err = sql.Open("mysql", "root:12306A.12306A@tcp(localhost:3310)/12306a_yutianneng_test")
	Db.SetMaxOpenConns(0)
	if err != nil {
		panic(err.Error())
	}
	CreateAllTables()
}

func DropTable(tableName string) {
	sqlStr := "drop table " + tableName
	Db.Exec(sqlStr)
}
func DropAllTables() {
	DropTable("train_pool")
	DropTable("total_train_no")

}

func CreateTable(sqlStr string) {
	Db.Exec(sqlStr)
}

var TICKETPOOLNUM = 10

func CreateTableTicketPool() {

	for i := 1; i <= TICKETPOOLNUM; i++ {
		sqlStr := "create table if not exists ticket_pool_" + strconv.Itoa(i) + "(" +
			"id int primary key auto_increment," +
			"train_no varchar(10)," +
			"initial_time datetime," +
			"start_time datetime," +
			"start_station varchar(50)," +
			"end_time datetime," +
			"end_station varchar(50)," +
			"carriage_no int," +
			"carriage_class varchar(10)," +
			"seat_no varchar(10)," +
			"ticket_price double" +
			");"
		_, err := Db.Exec(sqlStr)
		if err != nil {
			fmt.Println("create table ticket_pool failed, err:", err)
			return
		}
	}
}

func DropSplitTableTicketPool() {
	for i := 0; i < TICKETPOOLNUM; i++ {
		sqlStr := "drop table ticket_pool_" + strconv.Itoa(i) + ";"
		Db.Exec(sqlStr)
	}
}

//按车次对ticket_pool分表，原1900万条记录，分成100个表
//按照TrainNo分表，因为是按照车次去统计座位
//暂时不分表
//var TICKETPOOLNUM = 10
//
//func SplitTableTicketPoolByTrainNo() {
//
//	for i := 0; i < TICKETPOOLNUM; i++ {
//		sqlStr := "create table if not exists ticket_pool_" + strconv.Itoa(i) + "(" +
//			"id int primary key," +
//			"train_no varchar(10)," +
//			"initial_time datetime," +
//			"start_time datetime," +
//			"start_station varchar(50)," +
//			"end_time datetime," +
//			"end_station varchar(50)," +
//			"carriage_no int," +
//			"carriage_class varchar(10)," +
//			"seat_num int," +
//			"ticket_price double," +
//			"FOREIGN KEY() REFERENCES tb_dept1(id)" +
//			");"
//		Db.Exec(sqlStr)
//	}
//}
//
////暂时分为10个表作为测试
//var SEATPOOLNUM=10
//func SplitTableSeatPool()  {
//	for i:=0;i<SEATPOOLNUM;i++{
//		sqlStr:="create table if not exists seat_pool_"+strconv.Itoa(i)+" (" +
//			"id int primary key auto_increment," +
//			"ticket_pool_id int," +
//			"train_no varchar(10)," +//用于分表
//			"seat_no varchar(10)" +
//			");"
//		CreateTable(sqlStr)
//	}
//}

func CreateAllTables() {
	//创建station_provinve_city表
	sqlStr1 := "create table if not exists station_province_city (" +
		"id int primary key auto_increment," +
		"province varchar(50)," +
		"city varchar(50)," +
		"city_code varchar(50)," +
		"station_name varchar(50)," +
		"station_telecode varchar(4)," +
		"station_spell varchar(255)" +
		");"

	CreateTable(sqlStr1)

	//创建total_train_no，不用于读写，只用来生成票池
	sqlStr2 := "create table if not exists total_train_no (" +
		"id int primary key auto_increment," +
		"train_no varchar(10)," +
		"station_num int," +
		"initial_time datetime," +
		"terminal_time datetime," +
		"station_no int," +
		"station_name varchar(50)," +
		"city_name varchar(50)," +
		"arrive_time datetime," +
		"depart_time datetime," +
		"duration datetime," +
		"mileage int," +
		"price double);"
	CreateTable(sqlStr2)

	sqlStr3 := "create table if not exists train_pool(" +
		"id int primary key auto_increment," +
		"initial_time datetime," +
		"terminal_time datetime," +
		"train_no varchar(10)," +
		"start_city varchar(50)," +
		//"start_station varchar(50)," +
		"start_time datetime," +
		"end_city varchar(50)," +
		//"end_station varchar(50)," +
		"end_time datetime);"
	CreateTable(sqlStr3)
	//创建分表
	//CreateTableTicketPool()

}
