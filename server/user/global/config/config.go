/**
 * @Author fzh
 * @Date 2020/2/1
 */
package config

type Config struct {
	Server   Server   `ini:"server"`
	Database Database `ini:"database"`
}

type Server struct {
	Port    int `ini:"port"`
	RpcPort int `ini:"rpc_port"`
}

type Database struct {
	Ip       string `ini:"ip"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Database string `ini:"database"`
	Charset  string `ini:"charset"`
	LogMode  string `ini:"log_mode"`
}

var (
	Cfg *Config = new(Config)
)
