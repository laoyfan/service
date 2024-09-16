package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"net/http"
	"service/internal/config"
	"sync"
)

var (
	limit *limiter.Limiter
	once  sync.Once
)

func init() {
	once.Do(func() {
		limit = tollbooth.NewLimiter(config.AppConfig.Limit, nil)
	})
}

func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := tollbooth.LimitByRequest(limit, c.Writer, c.Request); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "服务繁忙，请稍后再试...",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
