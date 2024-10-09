package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"service/config"
	"service/logger"
	"service/middleware"
	"service/redis"
	"service/router"
	"service/translator"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化各个模块
	if err := initModules(); err != nil {
		fmt.Printf("初始化失败:%v\n", err)
		os.Exit(1)
	}

	// 关闭控制台颜色
	gin.DisableConsoleColor()
	// 设置模式
	gin.SetMode(config.AppConfig.Debug)

	// 开启gin实例
	r := gin.New()
	setupMiddleware(r)

	// HTTP配置
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        router.Route(r),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx, cancel := createContextWithTraceID()
	defer redis.Close(ctx) // 在服务关闭时断开 Redis 连接
	defer cancel()

	// 开启服务
	go func() {
		startServer(ctx, server)
	}()

	// 优雅Shutdown（或重启）服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "关闭服务...")

	// 关闭服务
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(ctx, "服务关闭原因:", zap.Error(err))
	}
	logger.Info(ctx, "服务退出")
}

// 初始化模块
func initModules() error {
	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("初始化配置异常: %v", err)
	}
	fmt.Println("配置初始化成功")
	if err := logger.InitLogger(); err != nil {
		return fmt.Errorf("初始化日志异常: %v", err)
	}
	fmt.Println("日志初始化成功")
	if err := redis.InitRedis(); err != nil {
		return fmt.Errorf("redis 初始化失败: %v", err)
	}
	fmt.Println("redis初始化成功")
	if err := translator.InitTranslator(); err != nil {
		return fmt.Errorf("验证器 翻译器 初始化失败: %v", err)
	}
	fmt.Println("翻译器初始化成功")
	middleware.InitLimiter()
	fmt.Println("限流器初始化成功")
	return nil
}

// 设置中间件
func setupMiddleware(r *gin.Engine) {
	r.Use(
		middleware.Cors(),    // 跨域处理
		middleware.Limiter(), // 限流处理
		middleware.Logger(),  // 日志处理
		middleware.Error(),   // 异常处理
	)
}

// 创建包含 Trace ID 的上下文
func createContextWithTraceID() (context.Context, context.CancelFunc) {
	baseCtx := context.Background()
	traceID := uuid.New().String()
	ctx := context.WithValue(baseCtx, uuid.New().String(), traceID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	return ctx, cancel
}

// 启动 HTTP 服务器
func startServer(ctx context.Context, server *http.Server) {
	logger.Info(ctx, fmt.Sprintf("服务开启:%d", config.AppConfig.Port))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(ctx, "listen: %s\n", zap.Error(err))
	}
}
