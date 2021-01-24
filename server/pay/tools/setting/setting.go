// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package setting

import (
	"flag"
	"github.com/go-ini/ini"
	"pay/tools/logging"
	"time"
)

type Server struct {
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RunMode      string
}

var ServerSetting = &Server{}

// 配置路径和配置文件名称
var configPath = flag.String("configPath", "./config/", "设置程序的配置文件路径")
var configName = flag.String("configName", "pay-config.ini", "设置配置文件的名称")

// Setup 载入配置文件信息
func Setup() {
	// 读取命令行信息
	flag.Parse()
	Cfg, err := ini.Load(*configPath + "/" + *configName)
	if err != nil {
		logging.Fatal("加载配置文件%s\\%s失败\n", configPath, configName)
	}
	loadConfig(Cfg, "server", ServerSetting)
	ServerSetting.ReadTimeout *= time.Second
	ServerSetting.WriteTimeout *= time.Second
}

func loadConfig(Cfg *ini.File, section string, data interface{}) {
	err := Cfg.Section(section).MapTo(data)
	if err != nil {
		logging.Fatal("加载%s数据出错: %v", section, err)
	}
}
