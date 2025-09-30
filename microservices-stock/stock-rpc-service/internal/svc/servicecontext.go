package svc

import (
	"fmt"
	"stock-rpc-service/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ServiceContext 服务上下文，用于注入配置和依赖
type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB         // MySQL 数据库连接
	Redis  *redis.Client    // Redis 客户端
	Mutex  *redsync.Redsync // Redis 分布式锁
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

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  redisClient,
		Mutex:  mutex,
	}
}
