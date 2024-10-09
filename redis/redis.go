package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"service/config"
	"service/logger"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

var redisClients sync.Map

func InitRedis() error {
	for _, instance := range config.AppConfig.Redis.Instances {
		for _, db := range instance.DBs {
			client := redis.NewClient(&redis.Options{
				Addr:         fmt.Sprintf("%s:%d", instance.Addr, instance.Port),
				Password:     instance.Password,
				DB:           db,
				PoolSize:     100,
				MinIdleConns: 10,

				DialTimeout:  10 * time.Second,
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
				MaxRetries:   5,

				Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
					netDialer := &net.Dialer{
						Timeout:   5 * time.Second,
						KeepAlive: 5 * time.Minute,
					}
					return netDialer.DialContext(ctx, network, addr)
				},
			})
			if err := client.Ping(context.Background()).Err(); err != nil {
				return fmt.Errorf("无法连接到 Redis instance %s, db %d: %w", instance.Name, db, err)
			}
			key := fmt.Sprintf("%s_%d", instance.Name, db)
			redisClients.Store(key, client)
		}
	}
	return nil
}

func GetRedisClient(name string, db int) (*redis.Client, error) {
	key := fmt.Sprintf("%s_%d", name, db)
	if client, ok := redisClients.Load(key); ok {
		return client.(*redis.Client), nil
	}
	return nil, errors.New(fmt.Sprintf("未配置 Redis %s and db %d", name, db))
}

func Close(ctx context.Context) {
	var wg sync.WaitGroup
	redisClients.Range(func(key, value interface{}) bool {
		if client, ok := value.(*redis.Client); ok {
			wg.Add(1)
			go func(key interface{}, client *redis.Client) {
				defer wg.Done()
				if err := client.Close(); err != nil {
					logger.Error(
						ctx,
						"关闭失败 Redis client",
						zap.Any("key", key),
						zap.Error(err))
				} else {
					logger.Info(
						ctx,
						"关闭成功 Redis client",
						zap.Any("key", key))
				}
			}(key, client)
		}
		return true
	})
	wg.Wait()
}
