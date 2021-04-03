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
	Host         string
	UserName     string
	PassWord     string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

var DataBase = &database{}

type server struct {
	Name     string
	Host     string
	HttpAddr string
	RPCAddr  string
	HttpPort int
}

var Server = &server{}

type redis struct {
	Host        string
	MaxIdle     int
	IdleTimeout time.Duration
}

var Redis = &redis{}

type kafka struct {
	Host string
}

var Kafka = &kafka{}

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


type rpcTarget struct {
	Order string
	TicketPool string
}

var RpcTarget = &rpcTarget{}


var configFile = flag.String("configFile", "config/ticket-config.ini", "设置配置文件")

func init() {
	flag.Parse()
	cfg, err := ini.Load(*configFile)
	if err != nil {
		logging.Fatal("Setting -- Load config fail:", err)
	}
	mapToStruct(cfg.Section("server"), Server)
	mapToStruct(cfg.Section("database"), DataBase)
	mapToStruct(cfg.Section("redis"), Redis)
	mapToStruct(cfg.Section("kafka"), Kafka)
	mapToStruct(cfg.Section("consul"), Consul)
	mapToStruct(cfg.Section("zipkin"), Zipkin)
	mapToStruct(cfg.Section("RPCTarget"), RpcTarget)
	Redis.IdleTimeout = Redis.IdleTimeout * time.Second
}

func mapToStruct(section *ini.Section, v interface{}) {
	err := section.MapTo(v)
	if err != nil {
		logging.Fatal("An error in section [%v] :", section.Name(), err)
	}
}
