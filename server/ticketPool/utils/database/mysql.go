// @Author LiuYong
// @Created at 2021-02-20
// @Modified at 2021-02-20
package database

import (
	"common/tools/logging"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ticketPool/utils/setting"
	"time"
)

var (
	DB *gorm.DB
)

func Setup() {

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", setting.DataBase.UserName, setting.DataBase.PassWord, setting.DataBase.DBHost, setting.DataBase.DBName, setting.DataBase.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Error("Fail to open db connect:", err)
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(setting.DataBase.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(setting.DataBase.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	stats, err := json.Marshal(sqlDB.Stats())
	logging.Info("Mysql Connection Pool stats:" + string(stats))

	DB = db
}

func Close() {
	sqlDB, _ := DB.DB()
	sqlDB.Close()
	logging.Info("connection pool is closed")
}
