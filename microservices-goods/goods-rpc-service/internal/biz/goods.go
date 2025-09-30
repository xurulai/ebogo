package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"goods-rpc-service/internal/svc"
	goodspb "goods-rpc-service/proto/goods"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GoodsBiz 商品业务逻辑
type GoodsBiz struct {
	svcCtx *svc.ServiceContext
}

// NewGoodsBiz 创建商品业务实例
func NewGoodsBiz(svcCtx *svc.ServiceContext) *GoodsBiz {
	return &GoodsBiz{
		svcCtx: svcCtx,
	}
}

// 商品模型
type Goods struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	GoodsId     int64     `gorm:"column:goods_id;uniqueIndex;not null"`
	CategoryId  int64     `gorm:"column:category_id;not null"`
	BrandName   string    `gorm:"column:brand_name;not null"`
	Code        string    `gorm:"column:code;uniqueIndex;not null"`
	Status      int8      `gorm:"column:status;not null"`
	Title       string    `gorm:"column:title;not null"`
	MarketPrice int64     `gorm:"column:market_price;not null"`
	Price       int64     `gorm:"column:price;not null"`
	Brief       string    `gorm:"column:brief;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

// TableName 定义表名
func (Goods) TableName() string {
	return "xx_goods_query"
}

// RoomGoods 直播间商品关联模型
type RoomGoods struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	RoomId    int64     `gorm:"column:room_id;not null"`
	GoodsId   int64     `gorm:"column:goods_id;not null"`
	Weight    int       `gorm:"column:weight;not null;default:0"`
	IsCurrent int8      `gorm:"column:is_current;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName 定义表名
func (RoomGoods) TableName() string {
	return "xx_room_goods"
}

var (
	localCache  = sync.Map{} // 本地缓存
	channelName = "cache_invalidation_channel"
)

// 定义一个结构体来表示缓存数据和版本号
type CacheDataWithVersion struct {
	Data    *goodspb.Goods `json:"data"`
	Version int64          `json:"version"`
}

// 定义一个结构体来表示频道发布的内容
type CacheUpdateMessage struct {
	CacheKey string         `json:"cache_key"`
	Data     *goodspb.Goods `json:"data"`
	Version  int64          `json:"version"`
}

// GetGoodsByRoom 根据直播间 ID 查询直播间绑定的所有商品信息
func (b *GoodsBiz) GetGoodsByRoom(ctx context.Context, roomId int64) (*goodspb.GetGoodsByRoomResponse, error) {
	// 1. 先去 xx_room_goods 表，根据 room_id 查询出所有的 goods_id
	var roomGoodsList []RoomGoods
	err := b.svcCtx.DB.WithContext(ctx).
		Model(&RoomGoods{}).
		Where("room_id = ?", roomId).
		Order("weight").
		Find(&roomGoodsList).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("查询直播间商品关联失败", zap.Error(err))
		return nil, fmt.Errorf("查询直播间商品关联失败: %w", err)
	}

	// 处理数据：拿出所有的商品 ID 和当前正在讲解的商品 ID
	var (
		currGoodsId int64
		idList      = make([]int64, 0, len(roomGoodsList))
	)

	for _, obj := range roomGoodsList {
		idList = append(idList, obj.GoodsId)
		if obj.IsCurrent == 1 {
			currGoodsId = obj.GoodsId
		}
	}

	if len(idList) == 0 {
		return &goodspb.GetGoodsByRoomResponse{
			CurrentGoodsId: currGoodsId,
			Data:           []*goodspb.GoodsInfo{},
		}, nil
	}

	// 2. 根据商品ID列表查询商品详细信息
	var goodsList []Goods
	err = b.svcCtx.DB.WithContext(ctx).
		Model(&Goods{}).
		Where("goods_id in ?", idList).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{
				SQL:                "FIELD(goods_id,?)",
				Vars:               []interface{}{idList},
				WithoutParentheses: true,
			},
		}).
		Find(&goodsList).Error

	if err != nil {
		zap.L().Error("查询商品详情失败", zap.Error(err))
		return nil, fmt.Errorf("查询商品详情失败: %w", err)
	}

	// 拼装响应数据
	data := make([]*goodspb.GoodsInfo, 0, len(goodsList))
	for _, goods := range goodsList {
		data = append(data, &goodspb.GoodsInfo{
			GoodsId:     goods.GoodsId,
			CategoryId:  goods.CategoryId,
			Status:      int32(goods.Status),
			Title:       goods.Title,
			MarketPrice: fmt.Sprintf("%.2f", float64(goods.MarketPrice)/100),
			Price:       fmt.Sprintf("%.2f", float64(goods.Price)/100),
			Brief:       goods.Brief,
		})
	}

	return &goodspb.GetGoodsByRoomResponse{
		CurrentGoodsId: currGoodsId,
		Data:           data,
	}, nil
}

// GetGoodsDetail 根据商品ID获取商品详情（带缓存）
func (b *GoodsBiz) GetGoodsDetail(ctx context.Context, goodsId int64) (*goodspb.GetGoodsDetailResponse, error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("goods_detail_%d", goodsId)

	// 1. 首先尝试从本地缓存中获取数据
	if localCacheData, ok := localCache.Load(cacheKey); ok {
		currentData := localCacheData.(CacheDataWithVersion)
		// 检查本地缓存的版本号是否为最新
		latestVersion, err := b.getLatestVersionFromRedis(cacheKey)
		if err != nil {
			log.Printf("Failed to get latest version from Redis: %v", err)
			// 如果获取最新版本号失败，仍然返回本地缓存的数据
			log.Printf("Local cache hit (version may not be latest):%d, version: %d", goodsId, currentData.Version)
			return &goodspb.GetGoodsDetailResponse{
				Goods: currentData.Data,
			}, nil
		}
		if currentData.Version == latestVersion {
			log.Printf("Local cache hit:%d, version: %d", goodsId, currentData.Version)
			return &goodspb.GetGoodsDetailResponse{
				Goods: currentData.Data,
			}, nil
		} else {
			log.Printf("Local cache version for %d is outdated (local: %d, latest: %d), fetching from Redis", goodsId, currentData.Version, latestVersion)
			return b.getFromRedisAndUpdateLocalCache(ctx, cacheKey, goodsId)
		}
	}

	// 2. 尝试从 Redis 缓存中获取数据
	return b.getFromRedisAndUpdateLocalCache(ctx, cacheKey, goodsId)
}

// getFromRedisAndUpdateLocalCache 从 Redis 获取数据并更新本地缓存
func (b *GoodsBiz) getFromRedisAndUpdateLocalCache(ctx context.Context, cacheKey string, goodsId int64) (*goodspb.GetGoodsDetailResponse, error) {
	cachedData, err := b.svcCtx.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// 缓存命中
		var cacheData CacheDataWithVersion
		if err := json.Unmarshal([]byte(cachedData), &cacheData); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			return nil, fmt.Errorf("缓存数据反序列化失败")
		}
		// 更新本地缓存
		b.updateLocalCache(cacheKey, cacheData.Data, cacheData.Version)
		log.Printf("Cache hit for GoodsId: %s, version: %d", cacheKey, cacheData.Version)
		return &goodspb.GetGoodsDetailResponse{
			Goods: cacheData.Data,
		}, nil
	} else if err != nil {
		log.Printf("Failed to get data from cache: %v", err)
	} else {
		log.Printf("Cache miss for GoodsId: %s", cacheKey)
	}

	// 缓存未命中，从数据库中查询数据
	mutexname := fmt.Sprintf("lock_goods_detail_%s", cacheKey)
	mutex := b.svcCtx.Mutex.NewMutex(mutexname)

	if err := mutex.Lock(); err != nil {
		return nil, fmt.Errorf("获取分布式锁失败: %w", err)
	}
	defer mutex.Unlock()

	var goodsDetail Goods
	err = b.svcCtx.DB.WithContext(ctx).
		Model(&Goods{}).
		Where("goods_id = ?", goodsId).
		First(&goodsDetail).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("商品不存在")
		}
		log.Printf("Failed to query goods detail: %v", err)
		return nil, fmt.Errorf("查询商品详情失败: %w", err)
	}

	// 检查商品详情数据是否有效
	if goodsDetail.GoodsId == 0 || goodsDetail.Title == "" || goodsDetail.Price == 0 {
		log.Printf("Invalid goods detail data: %+v", goodsDetail)
		return nil, fmt.Errorf("商品数据无效")
	}

	// 构造返回的响应数据
	goodsDetailResp := &goodspb.Goods{
		GoodsId:    goodsDetail.GoodsId,
		CategoryId: goodsDetail.CategoryId,
		Status:     int32(goodsDetail.Status),
		Title:      goodsDetail.Title,
		Code:       goodsDetail.Code,
		BrandName:  goodsDetail.BrandName,
		Brief:      goodsDetail.Brief,
	}

	// 格式化价格字段
	if goodsDetail.MarketPrice > 0 {
		goodsDetailResp.MarketPrice = fmt.Sprintf("%.2f", float64(goodsDetail.MarketPrice)/100)
	} else {
		goodsDetailResp.MarketPrice = "0.00"
	}

	if goodsDetail.Price > 0 {
		goodsDetailResp.Price = fmt.Sprintf("%.2f", float64(goodsDetail.Price)/100)
	} else {
		goodsDetailResp.Price = "0.00"
	}

	// 生成新的版本号
	version := time.Now().UnixNano()

	// 构造缓存数据
	cacheData := CacheDataWithVersion{
		Data:    goodsDetailResp,
		Version: version,
	}

	// 将数据序列化为 JSON
	cachedBytes, err := json.Marshal(cacheData)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return nil, fmt.Errorf("缓存数据序列化失败")
	}

	// 将序列化后的数据写入 Redis 缓存
	baseTTL := 10 * time.Minute
	randomTTL := time.Duration(rand.Intn(5*60)) * time.Second
	totalTTL := baseTTL + randomTTL
	_, err = b.svcCtx.Redis.Set(ctx, cacheKey, cachedBytes, totalTTL).Result()
	if err != nil {
		log.Printf("Failed to set data in cache: %v", err)
	}

	// 将数据存入本地缓存
	b.updateLocalCache(cacheKey, goodsDetailResp, version)

	// 构造通知消息并发布到 Redis 频道
	message := CacheUpdateMessage{
		CacheKey: cacheKey,
		Data:     goodsDetailResp,
		Version:  version,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal update message: %v", err)
	}
	// 发布消息到频道
	n, err := b.svcCtx.Redis.Publish(ctx, channelName, messageBytes).Result()
	if err != nil {
		log.Printf("Failed to publish message to channel: %v", err)
	} else {
		log.Printf("Published message to channel %s, number of subscribers: %d", channelName, n)
	}

	log.Printf("Returning goods detail response: %+v, version: %d", goodsDetailResp, version)
	return &goodspb.GetGoodsDetailResponse{
		Goods: goodsDetailResp,
	}, nil
}

// UpdateGoodsDetail 更新商品详情
func (b *GoodsBiz) UpdateGoodsDetail(ctx context.Context, goodsId int64, newPrice int64) error {
	// 1. 更新数据库
	result := b.svcCtx.DB.WithContext(ctx).
		Model(&Goods{}).
		Where("goods_id = ?", goodsId).
		Updates(map[string]interface{}{
			"price": newPrice,
		})

	if result.Error != nil {
		log.Printf("Failed to update goods detail: %v", result.Error)
		return fmt.Errorf("更新商品详情失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected for goodsId: %d", goodsId)
		return fmt.Errorf("商品不存在")
	}

	// 2. 删除缓存
	cacheKey := fmt.Sprintf("goods_detail_%d", goodsId)
	_, err := b.svcCtx.Redis.Del(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Failed to delete cache: %v", err)
		return fmt.Errorf("删除缓存失败: %w", err)
	}

	// 3. 删除本地缓存
	localCache.Delete(cacheKey)

	log.Printf("Cache deleted for GoodsId: %d", goodsId)
	return nil
}

// getLatestVersionFromRedis 从 Redis 获取最新的版本号
func (b *GoodsBiz) getLatestVersionFromRedis(cacheKey string) (int64, error) {
	cachedData, err := b.svcCtx.Redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		if cachedData == "" {
			log.Printf("No cached data found for key: %s", cacheKey)
			return 0, nil
		}
		log.Printf("Failed to get data from Redis: %v", err)
		return 0, err
	}

	var cacheData CacheDataWithVersion
	if err := json.Unmarshal([]byte(cachedData), &cacheData); err != nil {
		log.Printf("Failed to unmarshal cached data: %v", err)
		return 0, err
	}

	return cacheData.Version, nil
}

// updateLocalCache 更新本地缓存
func (b *GoodsBiz) updateLocalCache(cacheKey string, data *goodspb.Goods, version int64) {
	// 获取当前本地缓存中的数据和版本号
	current, ok := localCache.Load(cacheKey)
	if ok {
		currentData := current.(CacheDataWithVersion)
		// 如果本地缓存的版本号大于等于收到的版本号，说明本地缓存是新的，不需要更新
		if currentData.Version >= version {
			log.Printf("Local cache for key %s is already up to date", cacheKey)
			return
		}
	}

	// 更新本地缓存
	newCacheData := CacheDataWithVersion{
		Data:    data,
		Version: version,
	}
	localCache.Store(cacheKey, newCacheData)
	log.Printf("Updated local cache for key: %s, version: %d", cacheKey, version)
}
