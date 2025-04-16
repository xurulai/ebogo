package biz

import (
	"context"
	"seckil_service/dao/redis"
	"seckil_service/proto"
	"seckil_service/rpc"
	"strconv"
	"time"

	"go.uber.org/zap"
)

//设置定时任务或者由运维人员在活动开始前将库存预热
//筛选热点数据，提前将商品信息存入redis
// 刷新指定商品的库存到 Redis

func refreshStockForGoods(ctx context.Context, goodsID int64,userID int64) error{
	stock,err := rpc.StockCli.GetStock(ctx,&proto.GetStockReq{
		GoodsId: goodsID,
	})
	if err != nil {
        zap.L().Error("获取库存失败", zap.Error(err), zap.Int64("goodsID", goodsID))
        return err
    }
	// 将库存信息同步到 Redis
    // 将库存信息同步到 Redis
    stockKey := "seckill:stock:" + strconv.FormatInt(goodsID, 10) // 修改这里
    err = redis.GetClient().Set(ctx, stockKey, stock.Stock, time.Minute*5).Err()
    if err != nil {
        zap.L().Error("同步库存到Redis失败", zap.Error(err), zap.Int64("goodsID", goodsID))
        return err
    }

    zap.L().Info("库存预热完成", zap.Int64("goodsID", goodsID), zap.Int64("stock", stock.Stock))
    return nil
}