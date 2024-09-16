package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"service/internal/config"
	"service/util"
	"sync"
	"time"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		logger = NewLogger()
		zap.ReplaceGlobals(logger)
	})
}

func Info(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

func NewLogger() *zap.Logger {
	director := config.AppConfig.Zap.Director
	level := config.AppConfig.Zap.Level
	maxAge := config.AppConfig.Zap.MaxAge
	format := config.AppConfig.Zap.Format
	stackTraceKey := config.AppConfig.Zap.StackTraceKey
	encodeLevel := config.AppConfig.Zap.EncodeLevel
	prefix := config.AppConfig.Zap.Prefix
	logInConsole := config.AppConfig.Zap.LoginConsole
	showLine := config.AppConfig.Zap.ShowLine

	if ok, _ := util.PathExists(director); !ok {
		fmt.Println("创建日志文件夹", director)
		err := os.Mkdir(director, os.ModePerm)
		if err != nil {
			fmt.Println("创建日志文件夹失败", err)
		}
	}

	cores := make([]zapcore.Core, 0, 7)
	for zLevel := getLevel(level); zLevel <= zapcore.FatalLevel; zLevel++ {

		var (
			writer  zapcore.WriteSyncer
			eConfig zapcore.EncoderConfig
		)

		fileWriter, err := rotatelogs.New(
			path.Join(director, "%Y-%m-%d", zLevel.String()+".log"),
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(maxAge*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		if err != nil {
			fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
			continue
		}

		if logInConsole {
			writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
		} else {
			writer = zapcore.AddSync(fileWriter)
		}

		eConfig = zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: stackTraceKey,
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   levelEncoder(encodeLevel),
			EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(t.Format(prefix + "2006/01/02 - 15:04:05.000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}

		if format == "json" {
			cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(eConfig), writer, levelPriority(zLevel)))
		} else {
			cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(eConfig), writer, levelPriority(zLevel)))
		}
	}

	log := zap.New(zapcore.NewTee(cores...))
	if showLine {
		log = log.WithOptions(zap.AddCaller())
	}
	return log
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
		return zapcore.WarnLevel
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
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
