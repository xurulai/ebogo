package mq

import (
	"context"
	"seckil_service/model"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// SendSeckillRequest 发送秒杀请求到 RocketMQ
func SendSeckillRequest(ctx context.Context, userId int64,goodsId int64,quantity int) error {
	// 生成一个新的秒杀请求
	request := model.NewSeckillRequest(userId, goodsId, quantity)
	// 序列化请求为 JSON
	msgBody, err := jsoniter.Marshal(request)
	if err != nil {
		zap.L().Error("序列化秒杀请求失败", zap.Error(err), zap.Any("request", request))
		return err
	}

	// 创建消息
	msg := primitive.Message{
		Topic: "seckill-request-topic", // 替换为你的秒杀请求主题名称
		Body:  msgBody,
	}

	// 发送消息（异步）
	err = Producer.SendAsync(ctx, func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			// 如果发送失败，记录日志
			zap.L().Error("异步发送秒杀请求失败", zap.Error(err), zap.Any("request", request))
		} else {
			// 如果发送成功，记录日志
			zap.L().Info("秒杀请求发送成功", zap.Any("request", request))
		}
	}, &msg)

	if err != nil {
		// 如果发送失败，记录日志并返回错误
		zap.L().Error("发送秒杀请求到 RocketMQ 失败", zap.Error(err), zap.Any("request", request))
		return err
	}

	return nil
}
