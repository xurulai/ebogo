package config

import "github.com/zeromicro/go-zero/zrpc"

// Config 库存 RPC 服务配置结构体
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
}
