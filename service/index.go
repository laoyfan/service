package service

import (
	"context"
	"encoding/json"
	"service/constant"
	"service/logger"
	"service/model"
	"service/redis"

	"go.uber.org/zap"
)

type IndexService struct{}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (s *IndexService) SaveRedisData(ctx context.Context, value model.RedisData) (result map[string]interface{}) {
	result = map[string]interface{}{
		"msg":  "success",
		"code": 200,
	}

	client, err := redis.GetRedisClient("default", 0)
	if err != nil {
		errMsg := "redis client error"
		logger.Error(ctx, errMsg, zap.Error(err))
		result["msg"] = errMsg
		result["code"] = 500
		return result
	}
	data, err := json.Marshal(value)
	if err != nil {
		errMsg := "json error"
		logger.Error(ctx, errMsg, zap.Error(err))
		result["msg"] = errMsg
		result["code"] = 500
		return result
	}
	err = client.HSet(ctx, constant.RedisHash, value.Id, data).Err()
	if err != nil {
		errMsg := "redis save error"
		logger.Error(ctx, errMsg, zap.Error(err))
		result["msg"] = errMsg
		result["code"] = 500
		return result
	}
	return result
}

func (s *IndexService) GetRedisData(ctx context.Context, id string) model.RedisData {
	client, err := redis.GetRedisClient("default", 0)
	if err != nil {
		errMsg := "redis client error"
		logger.Error(ctx, errMsg, zap.Error(err))
		return model.RedisData{}
	}
	result, err := client.HGet(ctx, constant.RedisHash, id).Result()
	if err != nil {
		errMsg := "redis get error"
		logger.Error(ctx, errMsg, zap.Error(err))
		return model.RedisData{}
	}
	var redisData model.RedisData
	err = json.Unmarshal([]byte(result), &redisData)
	if err != nil {
		errMsg := "json error"
		logger.Error(ctx, errMsg, zap.Error(err))
		return model.RedisData{}
	}
	return redisData
}
