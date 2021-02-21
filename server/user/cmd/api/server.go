/**
 * @Author fzh
 * @Date 2020/2/1
 */
package api

import (
	"common/rpc_manage"
	"common/tools/logging"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
	"rpc/user/userpb"
	"sync"
	"time"
	"user/api/rpcapi"
	"user/global/config"
	"user/global/database"
	"user/router"
)

var (
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API Server",
		Example: "...",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				if err := run(); err != nil {
					logging.Fatal(err)
				}
				defer wg.Done()
			}()
			go func() {
				if err := runRpc(); err != nil {
					logging.Fatal(err)
				}
				defer wg.Done()
			}()
			wg.Wait()
		},
	}
)

func setup() {
	// 初始化日志模块
	logging.Setup()

	// 读取配置
	initConfig()

	// 初始化数据库连接
	initDatabase()
}

func run() error {
	r := router.InitRouter()

	// 启动
	err := r.Run(fmt.Sprintf(":%d", config.Cfg.Server.Port))

	return err
}

func runRpc() error {
	rpcServer := rpc_manage.NewGRPCServer()
	userpb.RegisterUserServiceServer(rpcServer, new(rpcapi.UserService))
	logging.Info("rpc服务启动中...")
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Cfg.Server.RpcPort))
	if err != nil {
		logging.Fatal("rpc服务启动失败", err)
	}

	if err = rpcServer.Serve(listen); err != nil {
		logging.Fatal("rpc服务启动失败", err)
	}
	return err
}

// 加载配置文件
func initConfig() {
	if err := ini.MapTo(config.Cfg, "./config/config.ini"); err != nil {
		logging.Error(err)
	} else {
		logging.Info("配置文件加载成功")
	}
	fmt.Printf("%#v\n", config.Cfg)
}

// 初始化数据库连接
// 对全局变量DB初始化
func initDatabase() {
	// 连接数据库
	cfg := config.Cfg.Database

	// DSN
	// 示例: "root:123456@tcp(127.0.0.1:3306)/12306A_test?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Ip, cfg.Port, cfg.Database, cfg.Charset)
	logging.Debug("DSN:", dsn)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("数据库连接失败")
	} else {
		log.Println("数据库连接成功")
	}

	// 设置连接池
	if sqlDB, err := db.DB(); err != nil {
		log.Fatal("数据库连接池配置失败")
	} else {
		// 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// 设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(100)
		// 设置连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// 全局变量-数据库连接
	database.DB = db
}
