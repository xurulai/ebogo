package biz

import (
	"context"
	"fmt"
	"seckil_service/dao/mysql"
	"seckil_service/dao/redis"
	"seckil_service/model"
	"time"

	"go.uber.org/zap"
)

// handleRedisGetFail 处理Redis获取库存失败的情况
func handleRedisGetFail(ctx context.Context, request model.SeckillRequest, goodsId int64, stockKey string) {
	// 失败操作记录数+1
	err := incrementFailCount(ctx, goodsId)
	if err != nil {
		zap.L().Error("记录失败操作失败", zap.Error(err))
		return
	}

	// 直接从数据库扣减库存
	// 从数据库扣减库存
	rowsAffected, err := mysql.DecrementStock(ctx, goodsId, request.Quantity)
	if err != nil || rowsAffected == 0 {
		zap.L().Error("数据库扣减库存失败", zap.Error(err))
		incrementFailCount(ctx, goodsId)
		return
	}

	// 将数据库数据同步到Redis缓存

	// 查询新的库存值
	newStock, err := mysql.GetStockByGoodsId(ctx, goodsId)
	if err != nil {
		zap.L().Error("查询数据库库存失败", zap.Error(err))
		return
	}
	err = redis.GetClient().Set(ctx, stockKey, newStock.StockNum, time.Minute*5).Err()
	if err != nil {
		zap.L().Error("更新Redis库存失败", zap.Error(err))
		return
	}
}

func incrementFailCount(ctx context.Context, goodsId int64) error {
    key := fmt.Sprintf("fail_count_%d", goodsId)
    // 使用Redis INCR命令增加失败计数
    failCount, err := redis.GetClient().Incr(ctx, key).Result()
    if err != nil {
        return fmt.Errorf("记录失败操作失败: %v", err)
    }

    // 检查失败次数是否达到阈值
    if failCount >= 5 {
        // 关闭商品的抢购标记
        err = closeSeckillFlag(ctx, goodsId)
        if err != nil {
            return fmt.Errorf("关闭抢购标记失败: %v", err)
        }

        // 启动库存校准
        go startInventoryCalibration(ctx, goodsId)
    }

    return nil
}

func closeSeckillFlag(ctx context.Context, goodsId int64) error {
    // 设置熔断器，拒绝所有对该商品的访问请求
    return redis.GetClient().Set(ctx, fmt.Sprintf("seckill:closed:%d", goodsId), "true", time.Hour).Err()
}

func isSeckillClosed(ctx context.Context, goodsId int64) (bool, error) {
    closed, err := redis.GetClient().Get(ctx, fmt.Sprintf("seckill:closed:%d", goodsId)).Result()
    if err != nil {
        return false, nil
    }
    return closed == "true", nil
}