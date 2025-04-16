package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"seckil_service/biz"
	"seckil_service/dao/redis"
	"seckil_service/errno"
	"seckil_service/model"
	"seckil_service/proto"
	"seckil_service/rpc"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
)

// SeckillRequestHandle 是处理秒杀请求的消费者回调函数。
// SeckillRequestHandle 是处理秒杀请求的消费者回调函数。
func SeckillRequestHandle(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	//解决消费重复
	//顺序消费
	//库存预热
	for _, msg := range msgs {
		var request model.SeckillRequest
		var err error

		// 解析消息体
		err = json.Unmarshal(msg.Body, &request)
		if err != nil {
			zap.L().Error("解析秒杀请求消息失败", zap.Error(err), zap.Any("messageBody", msg.Body))
			continue
		}

		// 检查是否已经处理过该请求
		goodsId := request.GoodsId//用于分布式锁商品ID
		requestID := request.ID // 假设 SeckillRequest 中有一个唯一的 ID 字段
		exists, err := redis.GetClient().Exists(ctx, "seckill:request:"+requestID).Result()
		if err != nil {
			zap.L().Error("检查秒杀请求幂等性失败", zap.Error(err), zap.String("requestID", requestID))
			continue
		}

		if exists > 0 {
			zap.L().Info("秒杀请求已处理，跳过重复消费", zap.String("requestID", requestID))
			continue
		}

		// 检查熔断器状态
		//熔断器打开说明该商品的失败操作次数达到阈值，正在进行库存校准，不能下单。
		circuitOpen, err := isCircuitOpen(ctx, goodsId)
		if err != nil || circuitOpen {
			zap.L().Warn("熔断器打开，拒绝秒杀请求", zap.Int64("goodsId", goodsId), zap.Error(err))
			continue
		}

		// 构造分布式锁的 key，粒度细化到商品ID
		mutexName := fmt.Sprintf("stock_mutex_%d", goodsId)
		mutex := redis.Rs.NewMutex(mutexName)

		// 尝试获取锁。
		if err := mutex.Lock(); err != nil {
			zap.L().Error("获取分布式锁失败",
				zap.String("mutexName", mutexName),
				zap.Error(err))
			return consumer.ConsumeRetryLater, errno.GetLockFaild
		}
		// 确保在函数结束时释放锁。
		defer mutex.Unlock() // 确保在函数结束时释放锁。


		// 从 Redis 中检查库存
        stockKey := "seckill:stock:" + strconv.FormatInt(request.GoodsId,10)
        stock, err := redis.GetClient().Get(ctx, stockKey).Int64()
        if err != nil {
            zap.L().Error("获取Redis库存失败", zap.Error(err), zap.Int64("goodsID", request.GoodsId))
            //获取库存失败，进入异常处理流程
            biz.handleRedisGetFail(ctx, request, goodsId, stockKey)
            continue
        }

        if stock - int64(request.Quantity) < 0 {
            zap.L().Warn("库存不足，秒杀失败", zap.Int64("goodsID", request.GoodsId))
            // 释放分布式锁
            redis.GetClient().Del(ctx, stockKey)
            continue
        }

		// 更新 Redis 中的库存信息
        updatedStock := stock - int64(request.Quantity)
        err = redis.GetClient().Set(ctx, stockKey, updatedStock, time.Minute*5).Err()
        if err != nil {
            zap.L().Error("更新Redis库存失败", zap.Error(err), zap.Int64("goodsID", request.GoodsId))
            // Redis扣减失败，进入异常处理流程
		    (ctx, request, goodsId, stockKey)
            continue
        }

		// 异步任务将Redis库存同步到数据库
		go func() {
			syncRedisToDB(ctx, goodsId, updatedStock)
		}()

		 // 释放分布式锁，允许其他实例查询库存
		 redis.GetClient().Del(ctx,stockKey)

		// 调用 RPC 创建订单
		maxRetries := 3
		for attempt := 0; attempt < maxRetries; attempt++ {
			_, err = rpc.OrderCli.CreateOrder(ctx, &proto.CreateOrderReq{
				GoodsId: request.GoodsId,
				Num:     int32(request.Quantity),
				UserId:  request.UserId,
			})

			if err == nil {
				break // 成功创建订单，退出循环
			}

			// 记录重试日志
			zap.L().Warn("OrderCli.CreateOrder failed, retrying...",
				zap.Error(err),
				zap.Int("attempt", attempt+1),
				zap.Any("request", request))
		}

		if err != nil {
			// 如果重试后仍然失败，记录错误并返回 Rollback 状态
			zap.L().Error("OrderCli.CreateOrder failed after retries",
				zap.Error(err),
				zap.Any("request", request))
			return consumer.ConsumeRetryLater, err
		}

		// 记录秒杀成功的日志
		zap.L().Info("秒杀成功", zap.Any("request", request))
		// 将请求ID标记为已处理，设置过期时间（例如5分钟）
		err = redis.GetClient().Set(ctx, "seckill:request:"+requestID, "1", time.Minute*1).Err()
		if err != nil {
			zap.L().Error("标记秒杀请求幂等性失败", zap.Error(err), zap.String("requestID", requestID))
		}

		// 启动看门狗机制，定期续期
		go watchdog(ctx, requestID)
	}
	return consumer.ConsumeSuccess, nil
}

// watchdog 定期续期幂等性检查的键值对
func watchdog(ctx context.Context, requestID string) {
	for {
		// 检查键是否存在
		exists, err := redis.GetClient().Exists(ctx, "seckill:request:"+requestID).Result()
		if err != nil {
			zap.L().Error("检查秒杀请求幂等性失败", zap.Error(err), zap.String("requestID", requestID))
			return
		}

		if exists == 0 {
			// 键不存在，停止看门狗
			return
		}

		// 续期键的过期时间
		err = redis.GetClient().Expire(ctx, "seckill:request:"+requestID, time.Minute*5).Err()
		if err != nil {
			zap.L().Error("续期秒杀请求幂等性失败", zap.Error(err), zap.String("requestID", requestID))
			return
		}

		// 等待一段时间后再次检查
		time.Sleep(time.Second * 10)
	}
}
