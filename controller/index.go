package controller

import (
	"service/logger"
	"service/request"
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

func (c *IndexController) GetDeviceInfo(ctx *gin.Context) {
	var req request.GetDeviceInfoReq
	if err := c.Valid(ctx, &req); err != nil {
		return
	}
	device, code := c.service.GetDeviceInfo(ctx, req)
	c.Success(ctx, map[string]interface{}{
		"code": code,
		"data": device,
	})
}

func (c *IndexController) GetCpcTaskEnv(ctx *gin.Context) {
	//var req request.GetTaskEnvReq
	//if err := c.Valid(r, &req); err != nil {
	//	return
	//}
	var req request.GetDeviceInfoReq
	if err := c.Valid(ctx, &req); err != nil {
		return
	}
	logger.Info(ctx, "测试")
	c.service.GetDeviceInfo(ctx, req)

	c.Success(ctx, map[string]interface{}{
		"test": 1,
	})
}
