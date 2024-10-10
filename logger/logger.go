package logger

import (
	"context"
	"fmt"
	"os"
	"path"
	"service/config"
	"service/util"
	"time"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() error {
	var err error
	logger, err = NewLogger()
	if err == nil {
		zap.ReplaceGlobals(logger)
	}
	return err
}

func logWithTraceID(ctx context.Context, level string, msg string, fields ...zap.Field) {
	traceID, _ := ctx.Value("TraceID").(string)
	switch level {
	case "info":
		logger.With(zap.Any("trace_id", traceID)).WithOptions(zap.AddCallerSkip(2)).Info(msg, fields...)
	case "error":
		logger.With(zap.Any("trace_id", traceID)).WithOptions(zap.AddCallerSkip(2)).Error(msg, fields...)
	case "fatal":
		logger.With(zap.Any("trace_id", traceID)).WithOptions(zap.AddCallerSkip(2)).Fatal(msg, fields...)
	}
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, "info", msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, "error", msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logWithTraceID(ctx, "fatal", msg, fields...)
}

func NewLogger() (*zap.Logger, error) {
	director := config.AppConfig.Zap.Director
	level := config.AppConfig.Zap.Level
	maxAge := config.AppConfig.Zap.MaxAge
	maxSize := config.AppConfig.Zap.MaxSize       // 每个日志文件的最大大小（MB）
	maxBackups := config.AppConfig.Zap.MaxBackups // 保留的最大备份数量
	format := config.AppConfig.Zap.Format
	stackTraceKey := config.AppConfig.Zap.StackTraceKey
	encodeLevel := config.AppConfig.Zap.EncodeLevel
	logInConsole := config.AppConfig.Zap.LoginConsole
	showLine := config.AppConfig.Zap.ShowLine

	if ok, _ := util.PathExists(director); !ok {
		fmt.Println("创建日志文件夹", director)
		err := os.Mkdir(director, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("创建日志文件夹失败: %w", err)
		}
	}

	// 获取当前日期，用于日志文件命名
	date := time.Now().Format("2006-01-02")
	logFileName := path.Join(director, fmt.Sprintf("%s.log", date))

	// 使用 Lumberjack 进行日志轮转
	lumberjackWriter := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,    // MB
		MaxBackups: maxBackups, // 保留的最大备份数量
		MaxAge:     maxAge,     // 天
		Compress:   true,       // 是否压缩备份
	}
	// 创建输出到控制台的 WriteSyncer
	var writer zapcore.WriteSyncer
	if logInConsole {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackWriter))
	} else {
		writer = zapcore.AddSync(lumberjackWriter)
	}

	eConfig := zapcore.EncoderConfig{
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

	// 创建核心
	var core zapcore.Core
	if format == "json" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(eConfig), writer, levelPriority(getLevel(level)))
	} else {
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(eConfig), writer, levelPriority(getLevel(level)))
	}

	log := zap.New(core)
	if showLine {
		log = log.WithOptions(zap.AddCaller())
	}
	return log, nil
}

// 获取配置对应level
func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func levelEncoder(encode string) zapcore.LevelEncoder {
	switch encode {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func levelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	return func(l zapcore.Level) bool {
		return l >= level
	}
}
