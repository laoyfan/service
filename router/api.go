package router

import (
	"service/controller"
	"service/middleware"
	"service/service"

	"github.com/gin-gonic/gin"
)

// Api api
func Api(r *gin.Engine) {

	api := r.Group("/api", middleware.Auth())

	// 实例化控制器
	indexController := controller.NewIndexController(service.NewIndexService())
	{
		api.POST("/saveRedisData", indexController.SaveRedisData)
		api.GET("/getRedisData", indexController.GetRedisData)
	}

}
