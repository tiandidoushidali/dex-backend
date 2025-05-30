package main

import (
	"dex/app/api/internal/config"
	"dex/app/api/internal/handler"
	"dex/app/api/internal/svc"
	"dex/app/data/utility"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 时间使用UTC
	time.Local = time.UTC
	// 初始化utility
	utility.Setup()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册中间件
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					httpx.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
						"code":    500,
						"message": fmt.Sprintf("%v", err),
					})
				}
			}()
			next(w, r)
		}
	})
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("111")
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
