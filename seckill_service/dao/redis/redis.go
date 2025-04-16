package redis

import (
	"context"
	"fmt"
	"seckil_service/config"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

// redsync -> https://github.com/go-redsync/redsync

var (
	rc *redis.Client
	Rs *redsync.Redsync
)

func Init(cfg *config.RedisConfig) error {
	rc = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // 密码
		DB:       cfg.DB,       // 数据�?
		PoolSize: cfg.PoolSize, // 连接池大�?
	})
	err := rc.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetClient 返回全局 Redis 客户端实例
func GetClient() *redis.Client {
	if rc == nil {
		panic("Redis client is not initialized")
	}
	return rc
}
