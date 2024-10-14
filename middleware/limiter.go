package middleware

import (
	"net/http"
	"service/config"
	"service/logger"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

var (
	once    sync.Once
	limiter *rate.Limiter
)

// InitLimiter 初始化限流器
func InitLimiter() {
	once.Do(func() {
		limiter = rate.NewLimiter(rate.Limit(config.AppConfig.Limit), 10)
	})
}

func Limiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查请求是否被限流
		if !limiter.Allow() {
			logger.Warn(ctx, "请求被限流",
				zap.String("url", ctx.Request.URL.Path),
				zap.String("client_ip", ctx.ClientIP()),
			)
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "服务繁忙，请稍后再试...",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
