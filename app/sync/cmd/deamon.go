package cmd

import (
	"context"
	"dex/app/sync/internal/service"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

var deamonCmd = &cobra.Command{
	Use:   "deamon",
	Short: "deamon commands",
	Long:  "deamon commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deamon called", args)
		var license, _ = cmd.Flags().GetString("license")
		fmt.Println("flag license:", license)
		//fmt.Println(cmd.Usage())

		ctx, cancel := context.WithCancel(context.Background())
		// 最佳实践：加入非缓冲通道，避免一直阻塞
		// 原因总结：Go 运行时避免阻塞自身调度器 甚至影响这个那个调度系统
		// “非阻塞式尝试发送” —— 如果不能立即写入，就直接放弃该信号（丢弃）。
		syncExit := make(chan os.Signal, 1)
		signal.Notify(syncExit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			s := service.New()
			s.Run(ctx)
		}()

		select {
		case <-syncExit:
			fmt.Println("deamon exit")
			cancel()
		}
		wg.Wait()
		fmt.Println("deamon finish")
	},
}

func init() {
	rootCmd.AddCommand(deamonCmd)
}
