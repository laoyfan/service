package middleware

import (
	"net/http"
	"service/config"

	"github.com/gin-gonic/gin"
)

// Cors 中间件处理跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置跨域响应头
		setHeaders(c)

		// 调试模式下放行所有请求
		if config.AppConfig.Debug == "debug" {
			c.Next()
			return
		}

		// 校验跨域请求
		if !allowOrigins(c) {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "校验跨域失败",
			})
			c.Abort()
			return
		}

		// OPTIONS 方法直接返回
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 处理请求
		c.Next()
	}
}

// setHeaders 设置允许跨域请求的响应头
func setHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Methods", config.AppConfig.Cors.AllowMethods)
	c.Header("Access-Control-Allow-Headers", config.AppConfig.Cors.AllowHeaders)
	c.Header("Access-Control-Expose-Headers", config.AppConfig.Cors.ExposeHeaders)
	c.Header("Access-Control-Allow-Credentials", config.AppConfig.Cors.AllowCredentials)
	c.Header("Access-Control-Max-Age", config.AppConfig.Cors.MaxAge)
}

// allowOrigins 校验请求的来源是否在允许的列表中
func allowOrigins(c *gin.Context) bool {
	origin := c.GetHeader("Origin")
	for _, allowedOrigin := range config.AppConfig.Cors.AllowOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
