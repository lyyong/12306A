/**
 * @Author fzh
 * @Date 2020/2/1
 */
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"user/cmd/api"
)

var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("用户服务模块")
	},
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
