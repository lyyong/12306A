/*
* @Author: 余添能
* @Date:   2021/3/3 8:27 下午
 */
package settings

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
type target struct {
	Addr string
}

var Server = &server{}
var DB = &database{}
var RedisDB = &redis{}
var Target = &target{}

// 配置路径
var configFile = flag.String("ConfigFile", "./config/search-config.ini", "设置配置文件")

//var configFile = flag.String("ConfigFile","/Users/yutianneng/go/src/12306A/server/search/config/search-config.ini","")

// Setup 载入配置文件信息
func Setup() {
	// 读取命令行信息
	flag.Parse()
	cfg, err := ini.Load(*configFile)
	if err != nil {
		logging.Fatal("Setting -- Load config fail:", err)
	}
	loadConfig(cfg, "server", Server)
	loadConfig(cfg, "database", DB)
	loadConfig(cfg, "redis", RedisDB)
	loadConfig(cfg, "target", Target)

	Server.ReadTimeout *= time.Second
	Server.WriteTimeout *= time.Second
	RedisDB.WriteTimeout *= time.Second
	RedisDB.ReadTimeout *= time.Second
	RedisDB.IdelTimeout *= time.Minute
}

func loadConfig(Cfg *ini.File, section string, data interface{}) {
	err := Cfg.Section(section).MapTo(data)
	if err != nil {
		logging.Error("加载%s数据出错: %v", section, err)
	}
}
