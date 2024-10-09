package middleware

import (
	"net/http"
	"service/config"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

var limit *limiter.Limiter

func InitLimiter() {
	limit = tollbooth.NewLimiter(config.AppConfig.Limit, nil)
	limit.SetTokenBucketExpirationTTL(1 * time.Second)
}

func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求是否超出限流
		if err := tollbooth.LimitByRequest(limit, c.Writer, c.Request); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "服务繁忙，请稍后再试...",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
