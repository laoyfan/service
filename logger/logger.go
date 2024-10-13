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
	logConfig := config.AppConfig.Zap
	if err := ensureLogDirectoryExists(logConfig.Director); err != nil {
		return nil, err
	}

	writer := getLogWriter(logConfig.Director, logConfig.MaxSize, logConfig.MaxBackups, logConfig.MaxAge)
	if logConfig.LoginConsole {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer)
	}

	encoderConfig := getEncoderConfig(logConfig.StackTraceKey, logConfig.EncodeLevel)
	var core zapcore.Core
	if logConfig.Format == "json" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, levelPriority(getLevel(logConfig.Level)))
	} else {
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writer, levelPriority(getLevel(logConfig.Level)))
	}

	log := zap.New(core)
	if logConfig.ShowLine {
		log = log.WithOptions(zap.AddCaller())
	}
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
	logFileName := path.Join(director, time.Now().Format("2006-01-02")+".log")
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
	})
}

// 获取编码配置
func getEncoderConfig(stackTraceKey, encodeLevel string) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: stackTraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   levelEncoder(encodeLevel),
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
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

// 日志级别编码器
func levelEncoder(encode string) zapcore.LevelEncoder {
	encoderMap := map[string]zapcore.LevelEncoder{
		"LowercaseLevelEncoder":      zapcore.LowercaseLevelEncoder,
		"LowercaseColorLevelEncoder": zapcore.LowercaseColorLevelEncoder,
		"CapitalLevelEncoder":        zapcore.CapitalLevelEncoder,
		"CapitalColorLevelEncoder":   zapcore.CapitalColorLevelEncoder,
	}
	if encoder, ok := encoderMap[encode]; ok {
		return encoder
	}
	return zapcore.LowercaseLevelEncoder
}
