// Package database
// @Author LiuYong
// @Created at 2021-02-04
package database

import (
	"common/tools/logging"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var client *gorm.DB
var db *sql.DB

func Setup(dbType, username, password, dbHost, dbname string) error {
	var err error
	client, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		dbHost,
		dbname)), &gorm.Config{})
	if err != nil {
		return err
	}
	db, _ = client.DB()
	db.SetConnMaxIdleTime(10)
	db.SetMaxOpenConns(100)
	client.Logger.LogMode(logger.Info)
	return nil
}

func Close() {
	if client != nil {
		db.Close()
	}
}

func Client() *gorm.DB {
	if client == nil {
		logging.Error("未使用database.Setup()进行数据库初始化")
		return nil
	}
	return client
}
