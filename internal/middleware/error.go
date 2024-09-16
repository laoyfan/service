package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"service/constant"
	"service/internal/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Error 中间件处理异常捕获
func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 处理断开连接情况
				if brokenPipe {
					logger.Error("请求连接断开",
						zap.String("url", c.Request.URL.Path),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.JSON(http.StatusOK, gin.H{
						"code": constant.ERROR,
						"msg":  "异常，请稍后重试",
						"data": nil,
					})
					c.Abort()
					return
				}

				// 记录异常日志和堆栈信息
				logger.Error("捕获到的错误",
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回服务器内部错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": constant.ERROR,
					"msg":  "服务器开小差，请稍后重试",
					"data": nil,
				})
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}
