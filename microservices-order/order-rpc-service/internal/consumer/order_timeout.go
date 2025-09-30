package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"order-rpc-service/internal/biz"
	"order-rpc-service/internal/svc"

	rocketmqConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OrderTimeoutConsumer 订单超时处理消费者
type OrderTimeoutConsumer struct {
	svcCtx *svc.ServiceContext
}

// NewOrderTimeoutConsumer 创建订单超时处理消费者
func NewOrderTimeoutConsumer(svcCtx *svc.ServiceContext) *OrderTimeoutConsumer {
	return &OrderTimeoutConsumer{
		svcCtx: svcCtx,
	}
}

// OrderTimeoutMessage 订单超时消息
type OrderTimeoutMessage struct {
	OrderId int64  `json:"orderId"`
	GoodsId int64  `json:"goodsId"`
	Num     int64  `json:"num"`
	Status  string `json:"status"`
}

// OrderTimeoutHandle 订单超时处理函数 - 保持原有逻辑
func (c *OrderTimeoutConsumer) OrderTimeoutHandle(ctx context.Context, msgs ...*primitive.MessageExt) (rocketmqConsumer.ConsumeResult, error) {
	for _, msg := range msgs {
		fmt.Printf("Received order timeout message: %s\n", string(msg.Body))

		var timeoutMsg OrderTimeoutMessage
		err := json.Unmarshal(msg.Body, &timeoutMsg)
		if err != nil {
			zap.L().Error("Failed to unmarshal timeout message", zap.Error(err))
			continue
		}

		// 检查订单状态
		var order biz.Order
		err = c.svcCtx.DB.WithContext(ctx).Where("order_id = ?", timeoutMsg.OrderId).First(&order).Error

		if err == gorm.ErrRecordNotFound {
			zap.L().Warn("Order not found for timeout processing", zap.Int64("orderId", timeoutMsg.OrderId))
			continue
		}

		if err != nil {
			zap.L().Error("Failed to query order", zap.Error(err))
			return rocketmqConsumer.ConsumeRetryLater, err
		}

		// 如果订单仍然是待支付状态，则进行超时处理
		if order.Status == "pending" {
			err = c.handleOrderTimeout(ctx, &timeoutMsg)
			if err != nil {
				zap.L().Error("Failed to handle order timeout", zap.Error(err))
				return rocketmqConsumer.ConsumeRetryLater, err
			}
		} else {
			zap.L().Info("Order status changed, skip timeout processing",
				zap.Int64("orderId", timeoutMsg.OrderId),
				zap.String("status", order.Status))
		}
	}

	return rocketmqConsumer.ConsumeSuccess, nil
}

// handleOrderTimeout 处理订单超时
func (c *OrderTimeoutConsumer) handleOrderTimeout(ctx context.Context, msg *OrderTimeoutMessage) error {
	// 1. 更新订单状态为超时
	err := c.svcCtx.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		err := tx.Model(&biz.Order{}).
			Where("order_id = ? AND status = ?", msg.OrderId, "pending").
			Update("status", "timeout").Error
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		// 更新订单详情状态
		err = tx.Model(&biz.OrderDetail{}).
			Where("order_id = ? AND status = ?", msg.OrderId, "pending").
			Update("status", "timeout").Error
		if err != nil {
			return fmt.Errorf("failed to update order detail status: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to update order timeout status: %w", err)
	}

	// 2. 回滚库存
	_, err = c.svcCtx.StockRpc.ReduceStock(ctx, &svc.ReduceStockInfo{
		GoodsId: msg.GoodsId,
		Num:     -msg.Num, // 负数表示回滚
		OrderId: msg.OrderId,
	})

	if err != nil {
		zap.L().Error("Failed to rollback stock for timeout order",
			zap.Error(err),
			zap.Int64("orderId", msg.OrderId))
		// 库存回滚失败不影响订单状态更新，但需要记录日志
	}

	zap.L().Info("Order timeout processed successfully",
		zap.Int64("orderId", msg.OrderId),
		zap.Int64("goodsId", msg.GoodsId),
		zap.Int64("num", msg.Num))

	return nil
}
