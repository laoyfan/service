package middleware

import (
	"net/http"
	"service/config"
	"service/logger"
	"sync"
	"time"

	"github.com/didip/tollbooth/limiter"

	"go.uber.org/zap"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
)

var (
	once  sync.Once
	limit *limiter.Limiter
)

// InitLimiter 初始化限流器
func InitLimiter() {
	once.Do(func() {
		limit = tollbooth.NewLimiter(config.AppConfig.Limit, nil)
		limit.SetTokenBucketExpirationTTL(1 * time.Second)
	})
}

func Limiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查请求是否超出限流
		if err := tollbooth.LimitByRequest(limit, ctx.Writer, ctx.Request); err != nil {
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
