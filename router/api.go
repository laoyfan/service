package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service/internal/logger"
)

// Api api
func Api(r *gin.Engine) {

	api := r.Group("/api")
	api.POST("/login", func(context *gin.Context) {
		logger.Info("测试")
		context.String(http.StatusOK, "hello World!")
	})

}
