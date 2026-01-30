// Package logger 统一日志记录
package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ruleback/internal/config"
)

var (
	globalLogger *zap.Logger
	globalSugar  *zap.SugaredLogger
	loggerOnce   sync.Once
	initErr      error
)

// Init 初始化日志（使用sync.Once确保只初始化一次）
func Init(cfg *config.LogConfig) error {
	loggerOnce.Do(func() {
		level := parseLevel(cfg.Level)

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		var encoder zapcore.Encoder
		if cfg.Format == "console" {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		}

		var writeSyncer zapcore.WriteSyncer
		if cfg.Output == "file" && cfg.FilePath != "" {
			file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				initErr = err
				return
			}
			writeSyncer = zapcore.AddSync(file)
		} else {
			writeSyncer = zapcore.AddSync(os.Stdout)
		}

		core := zapcore.NewCore(encoder, writeSyncer, level)
		globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		globalSugar = globalLogger.Sugar()
	})

	return initErr
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Field 创建日志字段
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// String 字符串字段
func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

// Int 整数字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 int64字段
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Uint 无符号整数字段
func Uint(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}

// Float64 浮点数字段
func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// Bool 布尔字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Err 错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Debug(msg, fields...)
	}
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Warn(msg, fields...)
	}
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

// Fatal 致命错误日志（会退出程序）
func Fatal(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Fatal(msg, fields...)
	}
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	if globalSugar != nil {
		globalSugar.Debugf(template, args...)
	}
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	if globalSugar != nil {
		globalSugar.Infof(template, args...)
	}
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	if globalSugar != nil {
		globalSugar.Warnf(template, args...)
	}
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	if globalSugar != nil {
		globalSugar.Errorf(template, args...)
	}
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	if globalSugar != nil {
		globalSugar.Fatalf(template, args...)
	}
}

// GetLogger 获取原始Logger实例
func GetLogger() *zap.Logger {
	return globalLogger
}

// GetSugar 获取Sugar实例
func GetSugar() *zap.SugaredLogger {
	return globalSugar
}

// Sync 同步日志缓冲区
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// WithFields 创建带预设字段的Logger
func WithFields(fields ...zap.Field) *zap.Logger {
	if globalLogger != nil {
		return globalLogger.With(fields...)
	}
	return nil
}
