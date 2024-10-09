package router

import (
	"service/controller"
	"service/middleware"
	"service/service"

	"github.com/gin-gonic/gin"
)

// Api api
func Api(r *gin.Engine) {

	api := r.Group("/api")

	// 实例化控制器
	indexController := controller.NewIndexController(service.NewIndexService())

	cpc := api.Group("/cpc", middleware.Auth(), middleware.Trace())
	{
		cpc.GET("/getInitDeviceInfo", indexController.GetDeviceInfo)
		cpc.POST("/getCpcTaskEnvV2", indexController.GetCpcTaskEnv)
	}

}
