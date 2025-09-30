package main

import (
	"flag"
	"fmt"

	"goods-rpc-service/internal/config"
	"goods-rpc-service/internal/server"
	"goods-rpc-service/internal/svc"
	goodspb "goods-rpc-service/proto/goods"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 启动参数：-f 指定配置文件路径
var configFile = flag.String("f", "etc/goods-rpc.yaml", "the config file")

// 程序入口：
// 1) 加载配置并初始化依赖
// 2) 创建并启动 gRPC 服务，注册 GoodsService 服务实现
// 3) 在开发/测试模式下开启 gRPC 反射便于调试
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		goodspb.RegisterGoodsServiceServer(grpcServer, server.NewGoodsServer(ctx))

		if c.RpcServerConf.Mode == service.DevMode || c.RpcServerConf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting gRPC server at %s...\n", c.RpcServerConf.ListenOn)
	s.Start()
}
