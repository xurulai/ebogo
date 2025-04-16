package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"seckil_service/config"
	"seckil_service/dao/mq"
	"seckil_service/handler"
	"seckil_service/logger"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
    var cfn string
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

    // 7. 初始化rocketmq
    err = mq.Init()
    if err != nil {
        panic(err)
    }

    // 初始化消费者
    c, err := rocketmq.NewPushConsumer(
        consumer.WithGroupName("order_srv_1"),
        consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
    )
    if err != nil {
        panic(err)
    }

    // 订阅主题
    err = c.Subscribe("seckill-request-topic", consumer.MessageSelector{},handler.SeckillRequestHandle)
    if err != nil {
        panic(err)
    }

    // 启动消费者
    err = c.Start()
    if err != nil {
        panic(err)
    }

    // 停止消费者的信号
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
    <-signalChan

    // 停止消费者
    c.Shutdown()
}

// MyConsumerListener 自定义消费者监听器
type MyConsumerListener struct {
}

// ConsumeMessage 实现顺序消费逻辑
func (l *MyConsumerListener) ConsumeMessage(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	return handler.SeckillRequestHandle(ctx, msgs...)
}