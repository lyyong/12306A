// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package setting

import (
	"common/tools/logging"
	"github.com/go-ini/ini"
	"time"
)

type database struct {
	UserName string
	PassWord string
	DBName string
	Charset string
	MaxIdleConns int
	MaxOpenConns int
}

var DataBase = &database{}


type server struct {
	HttpAddr string
	RpcAddr string
}

var Server = &server{}

type redis struct {
	Host string
	MaxIdle int
	IdleTimeout time.Duration
}

var Redis = &redis{}


func init() {
	cfg, err := ini.Load("config/indent-config.ini")
	if err != nil {
		logging.Fatal("Setting -- Load config fail:", err)
	}
	mapToStruct(cfg.Section("server"), Server)
	mapToStruct(cfg.Section("database"), DataBase)
	mapToStruct(cfg.Section("redis"), Redis)

	Redis.IdleTimeout = Redis.IdleTimeout * time.Second
}

func mapToStruct(section *ini.Section, v interface{}) {
	err := section.MapTo(v)
	if err != nil {
		logging.Fatal("An error in section [%v] :", section.Name(), err)
	}
}