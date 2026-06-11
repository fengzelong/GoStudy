package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Init 初始化全局 zap 日志，dev 环境同时输出到控制台。
func Init(env string) (func(), error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
		}),
		zap.InfoLevel,
	)

	cores := []zapcore.Core{fileCore}
	if env != "prod" {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.DebugLevel,
		))
	}

	zap.ReplaceGlobals(zap.New(zapcore.NewTee(cores...), zap.AddCaller()))

	return func() {
		_ = zap.L().Sync()
	}, nil
}
