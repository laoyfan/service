package service

import (
	"context"
	"encoding/json"
	"service/constant"
	"service/redis"
	"service/request"
)

type IndexService struct{}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (s *IndexService) GetDeviceInfo(ctx context.Context, req request.GetDeviceInfoReq) (device map[string]interface{}, code int) {
	client, _ := redis.GetRedisClient("default", 10)

	deviceInfo, err := client.HGet(ctx, constant.DeviceHash, req.ProcessId).Result()
	if err != nil {
		return nil, 103
	}

	data := make(map[string]interface{})

	err = json.Unmarshal([]byte(deviceInfo), &data)
	if err != nil {
		return nil, 501
	}
	return data, 200
}
