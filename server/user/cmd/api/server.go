package api

import (
	"github.com/spf13/cobra"
	"user/routers"
)

var (
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API Server",
		Example: "...",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func setup() {
	// 读取配置

	// 初始化数据库连接

}

func run() error {
	r := routers.InitRouter()

	// TODO: 启动
	err := r.Run(":8888")

	return err
}
