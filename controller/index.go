package controller

import (
	"service/model"
	"service/service"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
	Controller
	service *service.IndexService
}

func NewIndexController(service *service.IndexService) *IndexController {
	return &IndexController{
		service: service,
	}
}

func (c *IndexController) SaveRedisData(ctx *gin.Context) {
	var req model.RedisData
	if err := c.Valid(ctx, &req); err != nil {
		return
	}
	result := c.service.SaveRedisData(ctx, req)
	c.Success(ctx, result)
}

func (c *IndexController) GetRedisData(ctx *gin.Context) {
	var req model.GetRedisData
	if err := c.Valid(ctx, &req); err != nil {
		return
	}
	result := c.service.GetRedisData(ctx, req.Id)

	c.Success(ctx, result)
}
