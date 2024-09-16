package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"service/internal/config"
	"service/internal/logger"
	"service/internal/middleware"
	"service/internal/redis"
	"service/router"
	"syscall"
	"time"
)

func init() {

}

func Server() {
	// 服务停止时断开redis连接
	defer redis.Close()
	// 关闭控制台颜色
	gin.DisableConsoleColor()
	// 设置模式
	gin.SetMode(gin.ReleaseMode)
	// 开启gin实例
	r := gin.New()
	// 全局处理中间件
	r.Use(
		middleware.Logger(),  //日志处理
		middleware.Cors(),    //跨域处理
		middleware.Error(),   //异常处理
		middleware.Limiter(), //限流处理
	)

	// HTTP配置
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        router.Route(r),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 开启服务
	go func() {
		logger.Info(fmt.Sprintf("服务开启:%d", config.AppConfig.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()
	// 优雅Shutdown（或重启）服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("服务关闭原因:", zap.Error(err))
	}
	logger.Info("服务退出")
}
