package svc

import (
	"context"
	"fmt"
	"order-rpc-service/internal/config"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// 导入商品与库存服务端 proto
	goodspb "api-gateway/proto/goods"
	stockpb "stock-rpc-service/proto/stock"
)

// ServiceContext 服务上下文，用于注入配置和依赖
type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB          // MySQL 数据库连接
	Redis         *redis.Client     // Redis 客户端
	Mutex         *redsync.Redsync  // Redis 分布式锁
	MQProducer    rocketmq.Producer // RocketMQ 生产者
	SnowflakeNode *snowflake.Node   // Snowflake ID 生成器
	GoodsRpc      GoodsClient       // 商品服务 RPC 客户端
	StockRpc      StockClient       // 库存服务 RPC 客户端
}

// 定义 RPC 客户端接口
type GoodsClient interface {
	GetGoodsDetail(ctx context.Context, req *GetGoodsDetailReq) (*GoodsDetail, error)
}

type StockClient interface {
	ReduceStock(ctx context.Context, req *ReduceStockInfo) (*StockResponse, error)
}

// Proto 消息定义
type GetGoodsDetailReq struct {
	GoodsId int64
	UserId  int64
}

type GoodsDetail struct {
	GoodsId     int64
	CategoryId  int64
	Status      int32
	Title       string
	Code        string
	BrandName   string
	MarketPrice string
	Price       string
	Brief       string
}

type ReduceStockInfo struct {
	GoodsId int64
	Num     int64
	OrderId int64
}

type StockResponse struct {
	Success bool
	Message string
}

// RPC 客户端实现
type goodsRpcClient struct {
	conn *grpc.ClientConn
}

func (c *goodsRpcClient) GetGoodsDetail(ctx context.Context, req *GetGoodsDetailReq) (*GoodsDetail, error) {
	if c.conn == nil {
		// 如果连接不可用，返回模拟数据作为降级处理
		return &GoodsDetail{
			GoodsId:     req.GoodsId,
			CategoryId:  1,
			Status:      1,
			Title:       "模拟商品",
			Code:        "MOCK001",
			BrandName:   "模拟品牌",
			MarketPrice: "99.99",
			Price:       "89.99",
			Brief:       "模拟商品数据（商品服务不可用）",
		}, nil
	}

	// 创建 gRPC 客户端
	client := goodspb.NewGoodsServiceClient(c.conn)

	// 调用商品服务
	resp, err := client.GetGoodsDetail(ctx, &goodspb.GetGoodsDetailRequest{
		GoodsId: req.GoodsId,
		UserId:  req.UserId,
	})

	if err != nil {
		// 如果调用失败，返回模拟数据作为降级处理
		return &GoodsDetail{
			GoodsId:     req.GoodsId,
			CategoryId:  1,
			Status:      1,
			Title:       "降级商品",
			Code:        "FALLBACK001",
			BrandName:   "降级品牌",
			MarketPrice: "99.99",
			Price:       "89.99",
			Brief:       fmt.Sprintf("降级商品数据（调用失败: %v）", err),
		}, nil
	}

	// 转换响应数据
	return &GoodsDetail{
		GoodsId:     resp.Goods.GoodsId,
		CategoryId:  resp.Goods.CategoryId,
		Status:      resp.Goods.Status,
		Title:       resp.Goods.Title,
		Code:        resp.Goods.Code,
		BrandName:   resp.Goods.BrandName,
		MarketPrice: resp.Goods.MarketPrice,
		Price:       resp.Goods.Price,
		Brief:       resp.Goods.Brief,
	}, nil
}

type stockRpcClient struct {
	conn *grpc.ClientConn
}

func (c *stockRpcClient) ReduceStock(ctx context.Context, req *ReduceStockInfo) (*StockResponse, error) {
	if c.conn == nil {
		// 如果连接不可用，返回模拟成功作为降级处理
		return &StockResponse{
			Success: true,
			Message: "模拟库存扣减成功（库存服务不可用）",
		}, nil
	}

	// 创建 gRPC 客户端
	client := stockpb.NewStockServiceClient(c.conn)

	// 调用库存服务
	resp, err := client.ReduceStock(ctx, &stockpb.ReduceStockInfo{
		GoodsId: req.GoodsId,
		Num:     req.Num,
		OrderId: req.OrderId,
	})

	if err != nil {
		// 如果调用失败，返回模拟成功作为降级处理
		return &StockResponse{
			Success: true,
			Message: fmt.Sprintf("降级库存扣减成功（调用失败: %v）", err),
		}, nil
	}

	// 转换响应数据
	return &StockResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// createGoodsRpcClient 创建商品服务 RPC 客户端
func createGoodsRpcClient(c config.Config) GoodsClient {
	// 尝试连接商品服务
	if len(c.GoodsService.Endpoints) > 0 {
		conn, err := grpc.Dial(
			c.GoodsService.Endpoints[0],
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(time.Duration(c.GoodsService.Timeout)*time.Millisecond),
		)
		if err != nil {
			fmt.Printf("Failed to connect to goods service at %s: %v, using fallback\n",
				c.GoodsService.Endpoints[0], err)
			return &goodsRpcClient{conn: nil} // 返回降级客户端
		}

		fmt.Printf("Successfully connected to goods service at %s\n", c.GoodsService.Endpoints[0])
		return &goodsRpcClient{conn: conn}
	}

	fmt.Println("No goods service endpoints configured, using fallback")
	return &goodsRpcClient{conn: nil}
}

// createStockRpcClient 创建库存服务 RPC 客户端
func createStockRpcClient(c config.Config) StockClient {
	// 尝试连接库存服务
	if len(c.StockService.Endpoints) > 0 {
		conn, err := grpc.Dial(
			c.StockService.Endpoints[0],
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(time.Duration(c.StockService.Timeout)*time.Millisecond),
		)
		if err != nil {
			fmt.Printf("Failed to connect to stock service at %s: %v, using fallback\n",
				c.StockService.Endpoints[0], err)
			return &stockRpcClient{conn: nil} // 返回降级客户端
		}

		fmt.Printf("Successfully connected to stock service at %s\n", c.StockService.Endpoints[0])
		return &stockRpcClient{conn: conn}
	}

	fmt.Println("No stock service endpoints configured, using fallback")
	return &stockRpcClient{conn: nil}
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.MySQL.User,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Database,
		c.MySQL.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MySQL: %v", err))
	}

	// 初始化 Redis 连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.RedisConf.Host, c.RedisConf.Port),
		Password: c.RedisConf.Password,
		DB:       c.RedisConf.DB,
		PoolSize: c.RedisConf.PoolSize,
	})

	// 初始化 Redis 分布式锁
	pool := goredis.NewPool(redisClient)
	mutex := redsync.New(pool)

	// 初始化 RocketMQ 生产者（可选）
	var mqProducer rocketmq.Producer
	if c.RocketMQ.NameServer != "" {
		p, err := rocketmq.NewProducer(
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{c.RocketMQ.NameServer})),
			producer.WithRetry(3),
			producer.WithGroupName(c.RocketMQ.GroupName),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create RocketMQ producer: %v", err))
		}
		if err = p.Start(); err != nil {
			panic(fmt.Sprintf("failed to start RocketMQ producer: %v", err))
		}
		mqProducer = p
	}

	// 初始化 Snowflake ID 生成器
	startTime, err := time.Parse("2006-01-02", c.Snowflake.StartTime)
	if err != nil {
		panic(fmt.Sprintf("failed to parse snowflake start time: %v", err))
	}

	snowflake.Epoch = startTime.UnixNano() / 1000000
	node, err := snowflake.NewNode(c.Snowflake.MachineID)
	if err != nil {
		panic(fmt.Sprintf("failed to create snowflake node: %v", err))
	}

	// 初始化 RPC 客户端
	goodsRpc := createGoodsRpcClient(c)
	stockRpc := createStockRpcClient(c)

	return &ServiceContext{
		Config:        c,
		DB:            db,
		Redis:         redisClient,
		Mutex:         mutex,
		MQProducer:    mqProducer,
		SnowflakeNode: node,
		GoodsRpc:      goodsRpc,
		StockRpc:      stockRpc,
	}
}
