package mysql

import (
	"context"
	"seckil_service/errno"
	"seckil_service/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DecrementStock 从数据库中扣减库存
func DecrementStock(ctx context.Context, goodsId int64, quantity int) (int64, error) {
	// 使用 gorm 的 WithContext 方法，将上下文传递给数据库操作
	result := db.WithContext(ctx).
		Model(&model.Goods{}).
		Where("goods_id = ? AND stock >= ?", goodsId, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	// 检查更新是否成功
	if result.Error != nil {
		zap.L().Error("数据库扣减库存失败", zap.Error(result.Error))
		return 0, result.Error
	}

	// 如果没有行被更新，返回错误
	if result.RowsAffected == 0 {
		zap.L().Error("库存不足或商品不存在", zap.Int64("goodsId", goodsId))
		return 0, errno.ErrStockInsufficient
	}

	// 返回影响的行数
	return result.RowsAffected, nil
}

// GetStockByGoodsId 根据商品 ID 查询库存信息。
func GetStockByGoodsId(ctx context.Context, goodsId int64) (*model.Stock, error) {
	// 初始化一个 Stock 结构体用于存储查询结果。
	var data model.Stock

	// 使用 GORM 查询库存记录。
	err := db.WithContext(ctx).
		Model(&model.Stock{}).          // 指定操作的模型为 Stock 表。
		Where("goods_id = ?", goodsId). // 根据商品 ID 查询。
		First(&data).                   // 获取第一条记录。
		Error                           // 获取查询结果的错误信息。

	// 如果查询失败且错误不是记录未找到，则返回 ErrQueryFailed。
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errno.ErrQueryFailed
	}

	// 记录查询结果的日志。
	zap.L().Info("查询到的库存信息", zap.Any("data", data))

	return &data, nil // 返回查询结果。
}