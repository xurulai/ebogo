package biz

import (
	"context"
	"fmt"
	"stock-rpc-service/internal/svc"
	stockpb "stock-rpc-service/proto/stock"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StockBiz 库存业务逻辑
type StockBiz struct {
	svcCtx *svc.ServiceContext
}

// NewStockBiz 创建库存业务实例
func NewStockBiz(svcCtx *svc.ServiceContext) *StockBiz {
	return &StockBiz{
		svcCtx: svcCtx,
	}
}

// Stock 库存模型
type Stock struct {
	ID       int64 `gorm:"primaryKey;autoIncrement"`
	GoodsId  int64 `gorm:"column:goods_id;uniqueIndex;not null"`
	StockNum int64 `gorm:"column:stocknum;not null;default:0"` // 库存数量
	Lock     int64 `gorm:"column:lock;not null;default:0"`     // 预扣存数
}

// TableName 声明表名
func (Stock) TableName() string {
	return "xx_stock"
}

// StockRecord 库存记录模型
type StockRecord struct {
	ID      int64 `gorm:"primaryKey;autoIncrement"`
	OrderId int64 `gorm:"column:order_id;not null"`
	GoodsId int64 `gorm:"column:goods_id;not null"`
	Num     int64 `gorm:"column:num;not null"`
	Status  int   `gorm:"column:status;not null;default:1"` // 1:已减少 3:已回滚
}

// TableName 声明表名
func (StockRecord) TableName() string {
	return "xx_stock_record"
}

// SetStock 设置库存
func (b *StockBiz) SetStock(ctx context.Context, goodsId, num int64) error {
	// 使用 GORM 的 WithContext 方法确保操作在指定的上下文中执行
	result := b.svcCtx.DB.WithContext(ctx).
		Model(&Stock{}).                // 指定操作的模型为 Stock 表
		Where("goods_id = ?", goodsId). // 根据商品 ID 查询库存记录
		FirstOrCreate(&Stock{           // 如果记录不存在则创建，否则获取第一条记录
			GoodsId:  goodsId,
			StockNum: num,
		})

	// 如果查询或创建失败，返回错误
	if result.Error != nil {
		return fmt.Errorf("设置库存失败: %w", result.Error)
	}

	// 如果没有影响任何行（即记录已存在且未更新），则手动更新库存数量
	if result.RowsAffected == 0 {
		return b.svcCtx.DB.WithContext(ctx).
			Model(&Stock{}).
			Where("goods_id = ?", goodsId).
			Update("stocknum", num).Error // 更新库存数量
	}

	return nil // 操作成功，返回 nil
}

// GetStockByGoodsId 根据商品 ID 查询库存信息
func (b *StockBiz) GetStockByGoodsId(ctx context.Context, goodsId int64) (*stockpb.GoodsStockInfo, error) {
	// 初始化一个 Stock 结构体用于存储查询结果
	var data Stock

	// 使用 GORM 查询库存记录
	err := b.svcCtx.DB.WithContext(ctx).
		Model(&Stock{}).                // 指定操作的模型为 Stock 表
		Where("goods_id = ?", goodsId). // 根据商品 ID 查询
		First(&data).                   // 获取第一条记录
		Error                           // 获取查询结果的错误信息

	// 如果查询失败且错误不是记录未找到，则返回错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询库存失败: %w", err)
	}

	// 记录查询结果的日志
	zap.L().Info("查询到的库存信息", zap.Any("data", data))

	// 将查询结果封装到 Protobuf 消息 GoodsStockInfo 中
	resp := &stockpb.GoodsStockInfo{
		GoodsId: data.GoodsId,  // 设置商品 ID
		Stock:   data.StockNum, // 设置库存数量
	}

	return resp, nil // 返回封装好的 Protobuf 消息和 nil 错误
}

// ReduceStock 减少库存，支持事务和分布式锁，确保操作的原子性
func (b *StockBiz) ReduceStock(ctx context.Context, goodsId, num, orderId int64) error {
	var data Stock

	// 构造分布式锁的 key
	mutexname := fmt.Sprintf("xx-stock-%d", goodsId)

	// 创建 Redis 分布式锁
	mutex := b.svcCtx.Mutex.NewMutex(mutexname)

	// 尝试获取锁
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("获取分布式锁失败: %w", err)
	}
	defer mutex.Unlock() // 确保在函数结束时释放锁

	// 使用 GORM 事务执行库存减少操作
	err := b.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 查询当前库存
		err := tx.WithContext(ctx).
			Model(&Stock{}).
			Where("goods_id = ?", goodsId).
			First(&data).Error
		if err != nil {
			zap.L().Error("查询库存失败", zap.Int64("goods_id", goodsId), zap.Error(err))
			return err
		}
		zap.L().Info("查询到的库存信息", zap.Any("data", data))

		// 计算可用库存（当前库存 - 锁定库存）
		availableStock := data.StockNum - data.Lock

		// 检查是否有足够的库存
		if availableStock < num {
			zap.L().Error("库存不足", zap.Int64("goods_id", goodsId), zap.Int64("available_stock", availableStock), zap.Int64("requested_num", num))
			return fmt.Errorf("库存不足")
		}

		// 减少库存并增加锁定库存
		data.StockNum -= num
		data.Lock += num

		// 更新库存记录
		err = tx.WithContext(ctx).
			Save(&data).Error
		if err != nil {
			zap.L().Error("更新库存失败", zap.Int64("goods_id", goodsId), zap.Error(err))
			return err
		}

		// 创建库存记录
		stockRecord := StockRecord{
			OrderId: orderId,
			GoodsId: goodsId,
			Num:     num,
			Status:  1, // 状态为 1 表示已减少
		}
		err = tx.WithContext(ctx).
			Model(&StockRecord{}).
			Create(&stockRecord).Error
		if err != nil {
			zap.L().Error("创建库存记录失败", zap.Error(err))
			return err
		}
		return nil
	})

	// 如果事务失败，返回错误
	if err != nil {
		zap.L().Error("减少库存失败", zap.Int64("goods_id", goodsId), zap.Error(err))
		return err
	}

	// 记录减少库存成功的日志
	zap.L().Info("减少库存成功",
		zap.Int64("goods_id", goodsId),
		zap.Int64("num", num),
		zap.Int64("new_stock_num", data.StockNum),
	)
	return nil
}

// RollbackStock 根据消息回滚库存，支持事务操作
func (b *StockBiz) RollbackStock(ctx context.Context, goodsId, rollbackNum, orderId int64) error {
	// 构造分布式锁的 key
	mutexName := fmt.Sprintf("xx-stock-%d", orderId)

	// 创建 Redis 分布式锁
	mutex := b.svcCtx.Mutex.NewMutex(mutexName)

	// 尝试获取锁
	if err := mutex.Lock(); err != nil {
		zap.L().Error("获取分布式锁失败",
			zap.String("mutexName", mutexName),
			zap.Error(err))
		return fmt.Errorf("获取分布式锁失败: %w", err)
	}
	// 确保在函数结束时释放锁
	defer mutex.Unlock()

	// 使用 GORM 事务执行库存回滚操作
	return b.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		var stockRecord StockRecord

		// 查询库存记录
		err := tx.WithContext(ctx).
			Model(&StockRecord{}).
			Where("order_id = ? and goods_id = ? and status = 1", orderId, goodsId).
			First(&stockRecord).Error

		// 如果记录不存在，则直接返回，不做处理
		if err == gorm.ErrRecordNotFound {
			zap.L().Warn("库存记录不存在，无需回滚",
				zap.Int64("orderId", orderId),
				zap.Int64("goodsId", goodsId))
			return nil
		}

		// 如果查询失败，记录错误并返回
		if err != nil {
			zap.L().Error("根据订单 ID 查询库存记录失败",
				zap.Error(err),
				zap.Int64("orderId", orderId),
				zap.Int64("goodsId", goodsId))
			return err
		}

		// 查询库存信息
		var stock Stock
		err = tx.WithContext(ctx).
			Model(&Stock{}).
			Where("goods_id = ?", goodsId).
			First(&stock).Error
		if err != nil {
			zap.L().Error("查询库存失败",
				zap.Error(err),
				zap.Int64("goodsId", goodsId))
			return err
		}

		// 回滚库存
		stock.StockNum += rollbackNum // 增加库存数量
		stock.Lock -= rollbackNum     // 减少锁定库存
		if stock.Lock < 0 {           // 如果锁定库存小于 0，表示回滚失败
			zap.L().Error("回滚库存失败，锁定库存不足",
				zap.Int64("goodsId", goodsId),
				zap.Int64("stockNum", stock.StockNum),
				zap.Int64("lock", stock.Lock))
			return fmt.Errorf("回滚库存失败，锁定库存不足")
		}

		// 更新库存记录
		err = tx.WithContext(ctx).Save(&stock).Error
		if err != nil {
			zap.L().Warn("回滚库存失败",
				zap.Int64("goodsId", stock.GoodsId),
				zap.Error(err))
			return err
		}

		// 更新库存记录状态为 3（表示已回滚）
		stockRecord.Status = 3
		err = tx.WithContext(ctx).Save(&stockRecord).Error
		if err != nil {
			zap.L().Warn("更新库存记录状态失败",
				zap.Int64("goodsId", stock.GoodsId),
				zap.Error(err))
			return err
		}

		return nil
	})
}
