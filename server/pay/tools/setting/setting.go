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

// 配置路径和配置文件名称
var configPath = flag.String("configPath", "./config/", "设置程序的配置文件路径")
var configName = flag.String("configName", "pay-config.ini", "设置配置文件的名称")

// Setup 载入配置文件信息
func Setup() {
	// 读取命令行信息
	flag.Parse()
	Cfg, err := ini.Load(*configPath + "/" + *configName)
	if err != nil {
		logging.Error("加载配置文件%s\\%s失败", configPath, configName)
	}
	loadConfig(Cfg, "server", Server)
	loadConfig(Cfg, "consul", Consul)
	loadConfig(Cfg, "zipkin", Zipkin)
	loadConfig(Cfg, "database", Database)
	loadConfig(Cfg, "redis", Redis)
	loadConfig(Cfg, "kafka", Kafka)
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
