package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"service/constant"
	"service/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Error 中间件处理异常捕获
func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

				// 处理断开连接情况
				if brokenPipe {
					logger.Error(ctx,
						"请求连接断开",
						zap.String("url", ctx.Request.URL.Path),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					ctx.JSON(http.StatusOK, gin.H{
						"code": constant.ERROR,
						"msg":  "异常，请稍后重试",
					})
					ctx.Abort()
					return
				}

				// 记录异常日志和堆栈信息
				logger.Error(ctx,
					"异常捕获",
					zap.String("type", "server_error"),
					zap.Any("error", err),
					zap.String("request", strings.ReplaceAll(string(httpRequest), "\r\n", " ")),
					zap.Strings("stack", formatStackTrace()),
				)

				// 返回服务器内部错误响应
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": constant.ERROR,
					"msg":  "服务器开小差，请稍后重试",
				})
				ctx.Abort()
			}
		}()

		// 继续处理请求
		ctx.Next()
	}
}

func formatStackTrace() []string {
	stack := string(debug.Stack())
	stackLines := strings.Split(stack, "\n")
	formattedStack := make([]string, 0, len(stackLines))
	for _, line := range stackLines {
		if strings.Contains(line, "goroutine") {
			continue // 跳过 goroutine 信息
		}
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 0 && strings.Contains(trimmedLine, "service") {
			formattedStack = append(formattedStack, trimmedLine)
		}
	}
	return formattedStack
}
