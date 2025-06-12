package goods

import (
	"context"
	"encoding/json"
	"fmt"
	"goods_srv/dao/mysql"
	"goods_srv/dao/redis"
	"goods_srv/errno"
	"goods_srv/proto"
	"log"
	"math/rand"
	"sync"
	"time"
)

// biz层业务代码
// biz -> dao

var (
	localCache  = sync.Map{} // 本地缓存
	// 定义一个频道名
	channelName = "cache_invalidation_channel"
)
// 定义一个结构体来表示缓存数据和版本号
type CacheDataWithVersion struct {
	Data     *proto.GoodsDetail `json:"data"`
	Version  int64              `json:"version"`
}
// 定义一个结构体来表示频道发布的内容
type CacheUpdateMessage struct {
	CacheKey  string              `json:"cache_key"`
	Data      *proto.GoodsDetail `json:"data"`
	Version   int64              `json:"version"`
}
// GetRoomGoodsListProto 根据直播间 ID 查询直播间绑定的所有商品信息，并组装成 protobuf 响应对象返回
func GetGoodsByRoom(ctx context.Context, roomId int64) (*proto.GoodsListResp, error) {
	// 1. 先去 xx_room_goods 表，根据 room_id 查询出所有的 goods_id
	objList, err := mysql.GetGoodsByRoomId(ctx, roomId)
	if err != nil {
		return nil, err // 如果查询失败，直接返回错误
	}

	// 处理数据
	// 1. 拿出所有的商品 ID
	// 2. 记住当前正在讲解的商品 ID
	var (
		currGoodsId int64                            // 当前正在讲解的商品 ID
		idList      = make([]int64, 0, len(objList)) // 存储所有商品 ID 的切片
	)

	// 遍历查询结果，提取商品 ID 和当前讲解的商品 ID
	for _, obj := range objList {
		fmt.Printf("obj:%#v\n", obj)         // 打印当前对象信息（调试用）
		idList = append(idList, obj.GoodsId) // 将商品 ID 添加到 idList 中
		if obj.IsCurrent == 1 {              // 如果当前对象是正在讲解的商品
			currGoodsId = obj.GoodsId // 记录当前正在讲解的商品 ID
		}
	}

	// 2. 再拿上面获取到的 goods_id 去 xx_goods 表查询所有的商品详细信息
	goodsList, err := mysql.GetGoodsByIdList(ctx, idList)
	if err != nil {
		return nil, err // 如果查询失败，直接返回错误
	}

	// 拼装响应数据
	data := make([]*proto.GoodsInfo, 0, len(goodsList)) // 创建一个存储商品信息的切片
	for _, goods := range goodsList {
		data = append(data, &proto.GoodsInfo{ // 创建一个 GoodsInfo 对象并添加到 data 切片中
			GoodsId:     goods.GoodsId,                                       // 商品 ID
			CategoryId:  goods.CategoryId,                                    // 商品分类 ID
			Status:      int32(goods.Status),                                 // 商品状态
			Title:       goods.Title,                                         // 商品标题
			MarketPrice: fmt.Sprintf("%.2f", float64(goods.MarketPrice/100)), // 商品市场价（单位转换为元）
			Price:       fmt.Sprintf("%.2f", float64(goods.Price/100)),       // 商品售价（单位转换为元）
			Brief:       goods.Brief,                                         // 商品简介
		})
	}

	// 创建并返回 protobuf 响应对象
	resp := &proto.GoodsListResp{
		CurrentGoodsId: currGoodsId, // 当前正在讲解的商品 ID
		Data:           data,        // 商品信息列表
	}
	return resp, nil
}

func GetGoodsDetailById(ctx context.Context, goodsId int64) (*proto.GoodsDetail, error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("goods_detail_%d", goodsId)
	// 1. 首先尝试从本地缓存中获取数据
	if localCacheData, ok := localCache.Load(cacheKey); ok {
		currentData := localCacheData.(CacheDataWithVersion)
		// 检查本地缓存的版本号是否为最新
		// 从 Redis 获取最新的版本号
		latestVersion, err := getLatestVersionFromRedis(cacheKey)
		if err != nil {
			log.Printf("Failed to get latest version from Redis: %v", err)
			// 如果获取最新版本号失败，仍然返回本地缓存的数据
			log.Printf("Local cache hit (version may not be latest):%d, version: %d", goodsId, currentData.Version)
			return currentData.Data, nil
		}
		if currentData.Version == latestVersion {
			log.Printf("Local cache hit:%d, version: %d", goodsId, currentData.Version)
			//本地缓存的是最新数据直接返回
			return currentData.Data, nil
		} else {
			log.Printf("Local cache version for %d is outdated (local: %d, latest: %d), fetching from Redis", goodsId, currentData.Version, latestVersion)
			// 本地缓存版本号不是最新的，从 Redis 获取最新数据
			return getFromRedisAndUpdateLocalCache(ctx, cacheKey,goodsId)
		}
	}
	
	// 2. 尝试从 Redis 缓存中获取数据
	return getFromRedisAndUpdateLocalCache(ctx, cacheKey,goodsId)
}


// 从 Redis 获取数据并更新本地缓存
func getFromRedisAndUpdateLocalCache(ctx context.Context, cacheKey string,goodsId int64) (*proto.GoodsDetail, error) {
	cachedData, err := redis.GetClient().Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// 缓存命中
		var cacheData CacheDataWithVersion
		// 将缓存中的 JSON 数据反序列化
		if err := json.Unmarshal([]byte(cachedData), &cacheData); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			return nil, errno.ErrQueryFailed
		}
		// 更新本地缓存
		updateLocalCache(cacheKey, cacheData.Data, cacheData.Version)
		log.Printf("Cache hit for GoodsId: %s, version: %d", cacheKey, cacheData.Version)
		return cacheData.Data, nil
	} else if err != nil {
		// 如果从 Redis 获取数据失败，记录日志
		log.Printf("Failed to get data from cache: %v", err)
	} else {
		// 缓存未命中
		log.Printf("Cache miss for GoodsId: %s", cacheKey)
	}

	// 缓存未命中，从数据库中查询数据
	// 构造分布式锁的 key。
	mutexname := fmt.Sprintf("lock_goods_detail_%s", cacheKey)

	// 创建 Redis 分布式锁。
	mutex := redis.Rs.NewMutex(mutexname)

	// 尝试获取锁。
	if err := mutex.Lock(); err != nil {
		return nil, errno.ErrGetLockFailed
	}
	defer mutex.Unlock() // 确保在函数结束时释放锁。

	goodsDetail, err := mysql.GetGoodsDetailById(ctx, goodsId)
	if err != nil {
		log.Printf("Failed to query goods detail: %v", err)
		return nil, errno.ErrQueryFailed
	}

	// 检查查询结果是否为空
	if goodsDetail == nil {
		log.Printf("Goods detail not found for GoodsId: %s", cacheKey)
		return nil, errno.ErrGoodsDetailNull
	}

	// 检查商品详情数据是否有效
	if goodsDetail.GoodsId == 0 || goodsDetail.Title == "" || goodsDetail.Price == 0 {
		log.Printf("Invalid goods detail data: %+v", goodsDetail)
		return nil, errno.ErrGoodsDetailNull
	}

	// 构造返回的响应数据
	resp := &proto.GoodsDetail{
		GoodsId:    goodsDetail.GoodsId,
		CategoryId: goodsDetail.CategoryId,
		Status:     int32(goodsDetail.Status),
		Title:      goodsDetail.Title,
		Code:       goodsDetail.Code,      // 商品编码
		BrandName:  goodsDetail.BrandName, // 商品品牌名称
		Brief:      goodsDetail.Brief,
	}

	// 格式化市场价格和价格字段
	if goodsDetail.MarketPrice > 0 {
		resp.MarketPrice = fmt.Sprintf("%.2f", float64(goodsDetail.MarketPrice)/100)
	} else {
		resp.MarketPrice = "0.00"
		log.Printf("MarketPrice is zero or invalid for GoodsId: %s", cacheKey)
	}

	if goodsDetail.Price > 0 {
		resp.Price = fmt.Sprintf("%.2f", float64(goodsDetail.Price)/100)
	} else {
		resp.Price = "0.00"
		log.Printf("Price is zero or invalid for GoodsId: %s", cacheKey)
	}

	// 生成新的版本号
	version := time.Now().UnixNano()

	// 构造缓存数据和版本号
	cacheData := CacheDataWithVersion{
		Data:    resp,
		Version: version,
	}

	// 将数据序列化为 JSON
	cachedBytes, err := json.Marshal(cacheData)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return nil, errno.ErrQueryFailed
	}

	// 将序列化后的数据写入 Redis 缓存
	baseTTL := 10 * time.Minute
	randomTTL := time.Duration(rand.Intn(5*60)) * time.Second
	totalTTL := baseTTL + randomTTL
	_, err = redis.GetClient().Set(ctx, cacheKey, cachedBytes, totalTTL).Result()
	if err != nil {
		log.Printf("Failed to set data in cache: %v", err)
	}

	// 将数据存入本地缓存
	updateLocalCache(cacheKey, resp, version)

	// 构造通知消息并发布到 Redis 频道
	message := CacheUpdateMessage{
		CacheKey: cacheKey,
		Data:     resp,
		Version:  version,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal update message: %v", err)
	}
	// 发布消息到频道
	n, err := redis.GetClient().Publish(ctx, channelName, messageBytes).Result()
	if err != nil {
		log.Printf("Failed to publish message to channel: %v", err)
	} else {
		log.Printf("Published message to channel %s, number of subscribers: %d", channelName, n)
	}

	// 返回商品详情响应
	log.Printf("Returning goods detail response: %+v, version: %d", resp, version)
	return resp, nil
}

// 定义一个函数来订阅频道，并在收到消息时更新本地缓存和 Redis 缓存
func startSubscriber(ctx context.Context){
	subscriber := redis.GetClient().Subscribe(ctx,channelName)
	defer subscriber.Close()
	for{
		message,err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}
		fmt.Printf("Received message: %s on channel: %s\n", message.Payload, message.Channel)

		//解析消息内容
		var updateMessage CacheUpdateMessage
		err = json.Unmarshal([]byte(message.Payload), &updateMessage)
		if err != nil {
			log.Printf("Error parsing update message: %v", err)
			continue
		}

		
		// 更新本地缓存
		updateLocalCache(updateMessage.CacheKey, updateMessage.Data, updateMessage.Version)
		
		// 更新 Redis 缓存
		updateRedisCache(ctx,updateMessage.CacheKey, updateMessage.Data)
	}
}

// 更新本地缓存
func updateLocalCache(cacheKey string, data *proto.GoodsDetail, version int64) {
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
// 更新 Redis 缓存
func updateRedisCache(ctx context.Context,cacheKey string, data *proto.GoodsDetail) {
	if data == nil {
		// 如果数据为 nil，从 Redis 中删除该键
		_, err := redis.GetClient().Del(ctx, cacheKey).Result()
		if err != nil {
			log.Printf("Failed to delete from Redis: %v", err)
		}
		log.Printf("Deleted from Redis for key: %s", cacheKey)
		return
	}

	// 将数据序列化为 JSON
	cachedBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return
	}

	// 设置缓存的基础过期时间和随机过期时间
	baseTTL := 10 * time.Minute
	randomTTL := time.Duration(rand.Intn(5*60)) * time.Second
	totalTTL := baseTTL + randomTTL

	// 将序列化后的数据写入 Redis 缓存
	_, err = redis.GetClient().Set(ctx, cacheKey, cachedBytes, totalTTL).Result()
	if err != nil {
		log.Printf("Failed to set data in Redis: %v", err)
		return
	}

	log.Printf("Updated Redis cache for key: %s", cacheKey)
}

// UpdateGoodsDetail 更新商品详情，并删除缓存
func UpdateGoodsDetail(ctx context.Context, goodsId int64, newPrice int64) (*proto.Response, error) {
	// 1. 更新数据库
	err := mysql.UpdateGoodsDetail(ctx, goodsId, newPrice)
	if err != nil {
		log.Printf("Failed to update goods detail: %v", err)
		return nil, errno.ErrUpdateFailed
	}

	// 2. 删除缓存
	cacheKey := fmt.Sprintf("goods_detail_%d", goodsId)
	_, err = redis.GetClient().Del(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Failed to delete cache: %v", err)
		return nil, errno.ErrCacheDeleteFailed
	}

	log.Printf("Cache deleted for GoodsId: %d", goodsId)
	return &proto.Response{}, nil
}

// 从 Redis 获取最新的版本号
func getLatestVersionFromRedis(cacheKey string) (int64, error) {
	cachedData, err := redis.GetClient().Get(context.Background(), cacheKey).Result()
	if err != nil {
		if cachedData == ""{
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

//设置本地缓存
func setLocalCache(key string,value interface{},ttl time.Duration){
	//设置本地缓存
	localCache.Store(key,value)

	// 启动一个 goroutine 来处理缓存过期
	go func() {
		time.Sleep(ttl)
		localCache.Delete(key)
	}()
}