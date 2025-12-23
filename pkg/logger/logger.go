package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger(logFile string, level string) {
	// 设置级别
	var logLevel zapcore.Level
	switch level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	// 配置 Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// Console Encoder
	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// File Encoder
	fileEncoderConfig := encoderConfig
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// 文件切割
	writer := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	core := zapcore.NewTee(
		// 终端：使用 console 格式
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
		// 文件：使用 json 格式
		zapcore.NewCore(fileEncoder, zapcore.AddSync(writer), logLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
}
