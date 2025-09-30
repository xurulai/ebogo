package config

import "github.com/zeromicro/go-zero/zrpc"

// Config 订单 RPC 服务配置结构体
type Config struct {
	// gRPC 服务端配置（端口、模式、限流、链路追踪等）
	zrpc.RpcServerConf

	// MySQL 数据库配置
	MySQL struct {
		Host     string // 数据库主机地址
		Port     int    // 数据库端口
		User     string // 数据库用户名
		Password string // 数据库密码
		Database string // 数据库名称
		Charset  string // 字符集
	}

	// Redis 配置
	RedisConf struct {
		Host     string // Redis 主机地址
		Port     int    // Redis 端口
		Password string // Redis 密码
		DB       int    // Redis 数据库编号
		PoolSize int    // 连接池大小
	}

	// RocketMQ 配置
	RocketMQ struct {
		NameServer string // RocketMQ 名称服务器地址
		GroupName  string // 生产者组名
		Topic      struct {
			CreateOrder             string // 创建订单主题
			PayTimeout              string // 支付超时主题
			StockRollback           string // 库存回滚主题
			CreateOrderSuccessfully string // 订单创建成功主题
		}
	}

	// Snowflake ID 生成器配置
	Snowflake struct {
		StartTime string // 开始时间
		MachineID int64  // 机器ID
	}

	// 外部服务配置
	GoodsService struct {
		Name      string   // 商品服务名称
		Endpoints []string // 商品服务端点列表
		Timeout   int64    // 超时时间（毫秒）
	}

	StockService struct {
		Name      string   // 库存服务名称
		Endpoints []string // 库存服务端点列表
		Timeout   int64    // 超时时间（毫秒）
	}

	// Consul 配置
	Consul struct {
		Addr string // Consul 地址
	}
}
