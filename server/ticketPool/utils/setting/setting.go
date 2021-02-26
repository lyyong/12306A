// @Author: KongLingWen
// @Created at 2021/2/6
// @Modified at 2021/2/6

package setting

import (
	"common/tools/logging"
	"flag"
	"github.com/go-ini/ini"
	"time"
)

type database struct {
	DBHost       string
	UserName     string
	PassWord     string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

var DataBase = &database{}

type server struct {
	Name	 string
	Host	 string
	RpcAddr  string
	RPCPort	 int
}

var Server = &server{}

type redis struct {
	Host        string
	MaxIdle     int
	IdleTimeout time.Duration
}

var Redis = &redis{}

type consul struct {
	Address     string
	Interval    int
	TTL         int
	ServiceHost string
	ServiceID   string
}
var Consul = &consul{}

type zipkin struct {
	ServiceID    string
	HttpEndpoint string
}
var Zipkin = &zipkin{}

var configFile = flag.String("configFile", "config/ticketPool-config.ini", "设置配置文件")

func InitSetting() {
	flag.Parse()
	cfg, err := ini.Load(*configFile)
	if err != nil {
		logging.Fatal("Setting -- Load config fail:", err)
	}
	mapToStruct(cfg.Section("server"), Server)
	mapToStruct(cfg.Section("database"), DataBase)
	mapToStruct(cfg.Section("redis"), Redis)
	mapToStruct(cfg.Section("consul"), Consul)
	mapToStruct(cfg.Section("zipkin"), Zipkin)
	Redis.IdleTimeout = Redis.IdleTimeout * time.Second
}

func mapToStruct(section *ini.Section, v interface{}) {
	err := section.MapTo(v)
	if err != nil {
		logging.Fatal("An error in section [%v] :", section.Name(), err)
	}
}
