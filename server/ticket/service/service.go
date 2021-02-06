package service

import (
	"common/tools/logging"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"ticket/utils/setting"
	"time"
)

var(
	db *gorm.DB
	redisPool *redis.Pool
)

func init(){
	logging.Info("init DataBase")

	var err error
	db, err = NewMysqlDB()
	if err != nil {
		logging.Error("Fail to init DB:", err)
	}

	redisPool = NewRedisPool()

	Close()
}

func NewMysqlDB() (*gorm.DB, error){

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=%s&parseTime=True&loc=Local", setting.DataBase.UserName, setting.DataBase.PassWord, setting.DataBase.DBName, setting.DataBase.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Error("Fail to open db connect:", err)
		return nil, err
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

	return db, nil
}

func NewRedisPool() *redis.Pool {
	return &redis.Pool {
		MaxIdle: setting.Redis.MaxIdle,
		IdleTimeout: setting.Redis.IdleTimeout,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", setting.Redis.Host)
			if err != nil {
				return nil, err
			}
			return conn, err
		},
	}
}

func Close(){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func(){
		<-c
		sqlDB,_ := db.DB()
		sqlDB.Close()
		redisPool.Close()
		logging.Info("connection pool colsed")
	}()
}