package cmd

import (
	"dex/app/sync/internal/config"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	userLicense string
	configPath  string
	rootCmd     = cobra.Command{
		Use:   "sync",
		Short: "执行同步操作",
		Long:  "执行同步操作明细",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("in command root", cmd, args)
			fmt.Println("configPath:", configPath)
			config.Setup(configPath)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "MIT", "user license default MIT")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Default().Fatal(err)
	}
}
