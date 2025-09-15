package initialize

import (
	"context"
	"fmt"
	"time"

	"userop_srv/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.Redis.Host, global.ServerConfig.Redis.Port),
		Password: "",  // 如果有密码，在这里设置
		DB:       0,   // 使用默认DB
		PoolSize: 100, // 连接池大小
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		zap.S().Errorf("Redis连接失败: %v", err)
		panic(err)
	}

	global.RedisClient = rdb
	zap.S().Info("Redis连接成功")
}
