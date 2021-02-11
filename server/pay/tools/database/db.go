// @Author LiuYong
// @Created at 2021-02-04
// @Modified at 2021-02-04
package database

import (
	"common/tools/logging"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var client *gorm.DB

func Setup(dbType, username, password, dbHost, dbname string) error {
	var err error
	client, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		dbHost,
		dbname))
	if err != nil {
		return err
	}
	client.DB().SetConnMaxIdleTime(10)
	client.DB().SetMaxOpenConns(100)
	client.LogMode(true)
	return nil
}

func Close() {
	if client != nil {
		_ = client.Close()
	}
}

func Client() *gorm.DB {
	if client == nil {
		logging.Error("未使用database.Setup()进行数据库初始化")
		return nil
	}
	return client
}
