package main

import (
	"context"
	"flag"
	"fmt"
	bloomfilter "goods_srv/BloomFilter"
	"goods_srv/config"
	"goods_srv/dao/mysql"
	"goods_srv/dao/redis"
	"goods_srv/handler"
	"goods_srv/logger"
	"goods_srv/proto"
	"goods_srv/registry"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GoodsServer 是一个 gRPC 服务结构体，实现了 proto.UnimplementedGoodsServer 接口。
type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func main() {
	var cfn string
	ctx := context.Background()
	// 0. 从命令行获取配置文件路径，默认值为 "./conf/config.yaml"
	// 例如：stock_service -conf="./conf/config_qa.yaml"
	flag.StringVar(&cfn, "conf", "./conf/config.yaml", "指定配置文件路径")
	flag.Parse()

	// 1. 加载配置文件
	err := config.Init(cfn)
	if err != nil {
		panic(err) // 如果加载配置文件失败，直接退出程序
	}

	// 2. 初始化日志模块
	err = logger.Init(config.Conf.LogConfig, config.Conf.Mode)
	if err != nil {
		panic(err) // 如果初始化日志模块失败，直接退出程序
	}

	// 3. 初始化 MySQL 数据库连接
	err = mysql.Init(config.Conf.MySQLConfig)
	if err != nil {
		panic(err) // 如果初始化 MySQL 数据库失败，直接退出程序
	}

	// 4.初始化 Redis 连接
	err = redis.Init(config.Conf.RedisConfig)
	if err != nil {
		panic(err) // 如果初始化 Redis 失败，直接退出程序
	}
	//5.初始化布隆过滤器并启动
	err = bloomfilter.InitBloomFilter(ctx)
	if err != nil {
		panic(err)
	}

	err = registry.Init(config.Conf.ConsulConfig.Addr)
	if err != nil {
		zap.L().Error("Failed to initialize Consul", zap.Error(err))
		// 可以选择退出或继续运行，取决于业务需求
		panic(err)
	}

	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Conf.Port))
	if err != nil {
		panic(err)
	}

	// 创建 gRPC 服务
	s := grpc.NewServer()
	// 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	// 注册股票服务到 gRPC 服务
	proto.RegisterGoodsServer(s, &handler.GoodsSrv{})
	// 启动 gRPC 服务
	go func() {
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	// 注册服务到 Consul
	err = registry.Reg.RegisterService(config.Conf.Name, config.Conf.IP, config.Conf.Port, nil)
	if err != nil {
		zap.L().Error("Failed to register service to Consul", zap.Error(err))
		// 可以选择退出或继续运行，取决于业务需求
		panic(err)

	}
	// 打印 gRPC 服务启动日志
	zap.L().Info(
		"rpc server start",
		zap.String("ip", config.Conf.IP),
		zap.Int("port", config.Conf.Port),
	)

	// 服务退出时注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // 等待退出信号

	// 注销服务
	serviceId := fmt.Sprintf("%s-%s-%d", config.Conf.Name, config.Conf.IP, config.Conf.Port)
	registry.Reg.Deregister(serviceId)
}
