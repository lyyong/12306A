/*
* @Author: 余添能
* @Date:   2021/1/23 1:10 上午
 */
package init_data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var err error

func init() {
	Db, err = sql.Open("mysql", "root:12345678@tcp(localhost:3306)/12306a_test")
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

	sqlStr3 := "create table if not exists train_pools(" +
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

	sqlStr4:="create table if not exists ticket_pools(" +
		"id int primary key auto_increment," +
		"train_no varchar(10)," +
		"start_time datetime," +
		"start_station varchar(50)," +
		"start_station_no int," +
		"end_time datetime," +
		"end_station varchar(50)," +
		"end_station_no int," +
		"carriage_no int," +
		"seat_class varchar(10)," +
		"seat_no int," +
		"ticket_price int" +
		");"
	CreateTable(sqlStr4)
	//创建分表
	//CreateTableTicketPool()

}
