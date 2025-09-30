package main

import (
	"flag"
	"fmt"

	"api-gateway/internal/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// 启动参数：-f 指定配置文件路径
var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

// 程序入口：
// 1) 加载配置
// 2) 创建 HTTP 服务器并注册路由
// 3) 启动服务，监听端口
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 全局请求日志中间件
	server.Use(rest.ToMiddleware(middleware.RequestLogMiddleware))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
