package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"service/internal/config"
	"sync"
	"time"
)

var (
	redisClients sync.Map
	once         sync.Once
)

func init() {
	once.Do(func() {
		for _, instance := range config.AppConfig.Redis.Instances {
			for _, db := range instance.DBs {
				client := redis.NewClient(&redis.Options{
					Addr:         fmt.Sprintf("%s:%d", instance.Addr, instance.Port),
					Password:     instance.Password,
					DB:           db,
					MinIdleConns: 10,
					Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
						netDialer := &net.Dialer{
							Timeout:   5 * time.Second,
							KeepAlive: 5 * time.Minute,
						}
						return netDialer.DialContext(ctx, network, addr)
					},
				})
				if err := client.Ping(context.Background()).Err(); err != nil {
					panic(fmt.Sprintf("无法连接到 Redis instance %s, db %d: %v", instance.Name, db, err))
				}
				key := fmt.Sprintf("%s_%d", instance.Name, db)
				redisClients.Store(key, client)
			}
		}
	})
}

func GetRedisClient(name string, db int) (*redis.Client, error) {
	key := fmt.Sprintf("%s_%d", name, db)
	if client, ok := redisClients.Load(key); ok {
		return client.(*redis.Client), nil
	}
	return nil, errors.New(fmt.Sprintf("配置错误 Redis %s and db %d", name, db))
}

func Close() {
	redisClients.Range(func(key, value interface{}) bool {
		if client, ok := value.(*redis.Client); ok {
			if err := client.Close(); err != nil {
				fmt.Printf("关闭失败 Redis client for %s: %v\n", key, err)
			} else {
				fmt.Printf("关闭成功 Redis client for %s closed\n", key)
			}
		}
		return true
	})
}
