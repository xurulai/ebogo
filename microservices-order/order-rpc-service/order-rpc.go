package main

import (
	"flag"
	"fmt"

	orderpb "microservices-order-proto/order"
	"order-rpc-service/internal/config"
	"order-rpc-service/internal/consumer"
	"order-rpc-service/internal/server"
	"order-rpc-service/internal/svc"

	"github.com/apache/rocketmq-client-go/v2"
	rocketmqConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order-rpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// 启动订单超时消费者（仅在 NameServer 配置存在时启用）
	if c.RocketMQ.NameServer != "" {
		timeoutConsumer := consumer.NewOrderTimeoutConsumer(ctx)
		mqConsumer, err := rocketmq.NewPushConsumer(
			rocketmqConsumer.WithGroupName("order_srv_timeout"),
			rocketmqConsumer.WithNsResolver(primitive.NewPassthroughResolver([]string{c.RocketMQ.NameServer})),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create timeout consumer: %v", err))
		}

		// 订阅订单超时主题
		err = mqConsumer.Subscribe(c.RocketMQ.Topic.PayTimeout, rocketmqConsumer.MessageSelector{}, timeoutConsumer.OrderTimeoutHandle)
		if err != nil {
			panic(fmt.Sprintf("failed to subscribe timeout topic: %v", err))
		}

		// 启动消费者
		err = mqConsumer.Start()
		if err != nil {
			panic(fmt.Sprintf("failed to start timeout consumer: %v", err))
		}
		defer mqConsumer.Shutdown()
	} else {
		fmt.Println("RocketMQ NameServer not configured, skipping consumer initialization")
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		orderpb.RegisterOrderServiceServer(grpcServer, server.NewOrderServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
