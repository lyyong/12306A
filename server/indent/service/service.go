package service

import (
	"common/tools/logging"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var(
	db *gorm.DB
	redisPool *redis.Pool
)

func init(){
	logging.Info("init DB")
	configPath := "config/ticket-config.ini"
	gormDB, err := NewMysqlDB(configPath)
	if err != nil {
		logging.Error("Fail to init DB:", err)
	}
	db = gormDB

	redisHost := ":6379"
	redisPool = NewRedisPool(redisHost)

	Close()
}

func NewMysqlDB(configPath string) (*gorm.DB, error){
	cfg, err := ini.Load(configPath)
	if err != nil {
		logging.Error("Fail to Load config file:", err)
		return nil, err
	}
	sec, err := cfg.GetSection("database")
	if err != nil {
		logging.Error("Fail to get section 'database':", err)
	}
	dbType := sec.Key("dbtype").String()
	userName := sec.Key("username").String()
	password := sec.Key("password").String()
	dbname := sec.Key("dbname").String()
	charset := sec.Key("charset").String()
	maxIdleConns, err := sec.Key("maxIdleConns").Int()
	maxOpenConns, err := sec.Key("MaxOpenConns").Int()
	if err != nil {
		logging.Error("error config:", err)
		return nil, err
	}
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=%s&parseTime=True&loc=Local", userName, password, dbname, charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Error("Fail to open db connect:", err)
		return nil, err
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	stats, err := json.Marshal(sqlDB.Stats())
	logging.Info(dbType + "Pool stats:" + string(stats))

	return db, nil
}

func NewRedisPool(server string) *redis.Pool {
	return &redis.Pool {
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", server)
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