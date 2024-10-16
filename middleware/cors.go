package middleware

import (
	"net/http"
	"service/config"

	"github.com/gin-gonic/gin"
)

var allowedOriginsMap map[string]struct{}

func InitAllowedOrigins(origins []string) {
	allowedOriginsMap = make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		allowedOriginsMap[origin] = struct{}{}
	}
}

// Cors 中间件处理跨域请求
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置跨域响应头
		setHeaders(ctx)

		// OPTIONS 方法直接返回
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 调试模式下放行所有请求
		if config.AppConfig.Debug == "debug" {
			ctx.Next()
			return
		}

		// 校验跨域请求
		if !allowOrigins(ctx) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "校验跨域失败",
			})
			ctx.Abort()
			return
		}

		// 处理请求
		ctx.Next()
	}
}

// setHeaders 设置允许跨域请求的响应头
func setHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
	ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type,New-Token,New-Expires-At")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Max-Age", "86400")
}

// allowOrigins 校验请求的来源是否在允许的列表中
func allowOrigins(ctx *gin.Context) bool {
	origin := ctx.GetHeader("Origin")
	_, exists := allowedOriginsMap[origin]
	return exists
}
