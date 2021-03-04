// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package setting

import (
	"common/tools/logging"
	"flag"
	"github.com/go-ini/ini"
	"time"
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

type kafka struct {
	Host string
}

var Server = &server{}
var Consul = &consul{}
var Zipkin = &zipkin{}
var Database = &database{}
var Redis = &redis{}
var Kafka = &kafka{}

// 配置路径
var configFile = flag.String("ConfigFile", "./config/search-config.ini", "设置配置文件")

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
