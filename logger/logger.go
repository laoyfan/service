package logger

import (
	"context"
	"fmt"
	"os"
	"path"
	"service/config"
	"service/util"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *zap.Logger
)

// InitLogger 初始化全局日志器
func InitLogger() error {
	var err error
	once.Do(func() {
		logger, err = NewLogger()
		if err == nil {
			zap.ReplaceGlobals(logger)
		}
	})
	return err
}

// NewLogger 创建新的日志器实例
func NewLogger() (*zap.Logger, error) {
	// 创建日志目录
	if err := ensureLogDirectoryExists(config.AppConfig.Zap.Director); err != nil {
		return nil, err
	}
	// 设置日志级别
	writer := getLogWriter(config.AppConfig.Zap.Director, config.AppConfig.Zap.MaxSize, config.AppConfig.Zap.MaxBackups, config.AppConfig.Zap.MaxAge)
	// debug模式输出控制台
	if config.AppConfig.Debug == "debug" {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer)
	}
	// 创建编码器配置
	encoderConfig := getEncoderConfig()
	var core zapcore.Core
	if config.AppConfig.Zap.Format == "json" {
		// 如果是JSON格式则使用JSONEncoder
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, levelPriority(getLevel(config.AppConfig.Zap.Level)))
	} else {
		// 如果是Console格式则使用ConsoleEncoder
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, levelPriority(getLevel(config.AppConfig.Zap.Level)))
	}

	log := zap.New(core)
	log = log.WithOptions(zap.AddCaller())

	return log, nil
}

// logWithTraceID 带有 TraceID 的日志记录
func logWithTraceID(ctx context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	traceID, _ := ctx.Value("TraceID").(string)
	if logger.Core().Enabled(level) {
		logger.With(zap.Any("trace_id", traceID)).WithOptions(zap.AddCallerSkip(2)).Log(level, msg, fields...)
	}
}

// Info 记录 Info 级别日志
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, zapcore.InfoLevel, msg, fields...)
}

// Warn 记录 Warn 级别日志
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, zapcore.WarnLevel, msg, fields...)
}

// Error 记录 Error 级别日志
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, zapcore.ErrorLevel, msg, fields...)
}

// Fatal 记录 Fatal 级别日志
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, zapcore.FatalLevel, msg, fields...)
}

// 创建日志目录
func ensureLogDirectoryExists(director string) error {
	if ok, _ := util.PathExists(director); !ok {
		fmt.Println("创建日志文件夹", director)
		if err := os.Mkdir(director, os.ModePerm); err != nil {
			return fmt.Errorf("创建日志文件夹失败: %w", err)
		}
	}
	return nil
}

// 获取日志文件写入器
func getLogWriter(director string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	logFileName := path.Join(director, "service.log")
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	})
}

// 获取编码配置
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// 获取日志级别
func getLevel(level string) zapcore.Level {
	levelMap := map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
	if zapLevel, exists := levelMap[level]; exists {
		return zapLevel
	}
	return zapcore.DebugLevel
}

// 获取日志级别优先级
func levelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	return func(l zapcore.Level) bool {
		return l >= level
	}
}
