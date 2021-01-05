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
		fmt.Println("欢迎使用")
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
