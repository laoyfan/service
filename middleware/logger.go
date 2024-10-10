package middleware

import (
	"service/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogLayout 日志layout结构体
type LogLayout struct {
	Time      time.Time `json:"time"`                // 请求的时间
	Status    int       `json:"status,omitempty"`    // HTTP响应状态码
	Method    string    `json:"method,omitempty"`    // HTTP方法
	Path      string    `json:"path,omitempty"`      // 请求的路径
	Query     string    `json:"query,omitempty"`     // 请求的Query参数
	IP        string    `json:"IP,omitempty"`        // 客户端IP
	UserAgent string    `json:"userAgent,omitempty"` // 用户代理
	Error     string    `json:"error,omitempty"`     // 错误信息
	Cost      float64   `json:"cost,omitempty"`      // 请求耗时
	Source    string    `json:"source,omitempty"`    // 请求来源
}

// Logger 返回一个 Gin 日志中间件
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()               // 记录请求开始时间
		path := ctx.Request.URL.Path      // 记录请求路径
		query := ctx.Request.URL.RawQuery // 记录请求Query参数

		ctx.Next()

		cost := time.Since(start)
		status := ctx.Writer.Status()

		// 构造日志信息
		layout := LogLayout{
			Time:      start,
			Status:    status,
			Method:    ctx.Request.Method,
			Path:      path,
			Query:     query,
			IP:        ctx.ClientIP(),
			UserAgent: ctx.Request.UserAgent(),
			Error:     ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			Cost:      cost.Seconds(),
			Source:    ctx.Request.Host,
		}

		// 根据响应状态码决定记录日志级别
		if status >= 400 {
			logger.Error(ctx, "请求错误", zap.Any("log", layout))
		} else {
			logger.Info(ctx, "请求成功", zap.Any("log", layout))
		}
	}
}
