// Package setting
// @Author liuYong
// @Created at 2020-01-05
package setting

import (
	"common/tools/logging"
	"flag"
	"time"

	"github.com/go-ini/ini"
)

type server struct {
	Name         string
	Host         string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RunMode      string
}

type consul struct {
	Address     string
	Interval    int
	TTL         int
	ServiceHost string
	ServiceID   string
}

type zipkin struct {
	ServiceID    string
	HttpEndpoint string
}

type database struct {
	Type     string
	Username string
	Password string
	Host     string
	DbName   string
}

type redis struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
	IdelTimeout  time.Duration
}

type rpcTarget struct {
	Order  string
	Ticket string
}

type kafka struct {
	Host string
}

var Server = &server{}
var Consul = &consul{}
var Zipkin = &zipkin{}
var Database = &database{}
var Kafka = &kafka{}
var RPCTarget = &rpcTarget{}
var Redis = &redis{}

// 配置路径
var configFile = flag.String("ConfigFile", "./config/candidate-config.ini", "设置配置文件")

// Setup 载入配置文件信息
func Setup() {
	// 读取命令行信息
	flag.Parse()
	cfg, err := ini.Load(*configFile)
	if err != nil {
		logging.Fatal("Setting -- Load config fail:", err)
	}
	loadConfig(cfg, "server", Server)
	loadConfig(cfg, "consul", Consul)
	loadConfig(cfg, "zipkin", Zipkin)
	loadConfig(cfg, "database", Database)
	loadConfig(cfg, "redis", Redis)
	loadConfig(cfg, "kafka", Kafka)
	loadConfig(cfg, "RPCTarget", RPCTarget)
	Server.ReadTimeout *= time.Second
	Server.WriteTimeout *= time.Second
	Redis.WriteTimeout *= time.Second
	Redis.ReadTimeout *= time.Second
	Redis.IdelTimeout *= time.Minute
}

func loadConfig(Cfg *ini.File, section string, data interface{}) {
	err := Cfg.Section(section).MapTo(data)
	if err != nil {
		logging.Error("加载%s数据出错: %v", section, err)
	}
}
