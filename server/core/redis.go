package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-makeadmin/config"
	"log"
	"sync"
	"time"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
	redisErr    error
)

func GetRedis() *redis.Client {
	redisOnce.Do(func() {
		redisClient, redisErr = initRedis()
	})
	if redisErr != nil {
		log.Fatal("GetRedis initRedis err: ", redisErr)
	}
	return redisClient
}

func CloseRedis() {
	if redisClient == nil {
		return
	}
	_ = redisClient.Close()
}

// initRedis 初始化redis客户端
func initRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL(config.Config.RedisUrl)
	if err != nil {
		return nil, err
	}
	opt.PoolSize = config.Config.RedisPoolSize
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
