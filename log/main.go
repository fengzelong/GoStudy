package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"time"
)

var infoLevel = "info level log"
var debugLevel = "debug level log"

var Logger *zap.Logger

func main() {
	// 6 初始化logger
	GetLogger()
	//Logger.Info(infoLevel, zap.Time("write info time", time.Now()))
	//Logger.Debug(debugLevel, zap.Time("debug time", time.Now()))

	r := gin.New()
	r.Use(GinLogger())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "v1.1",
		})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	r.Run(":8080")
}

// GinLogger 集成gin中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GetLogger 获取zap log对象
func GetLogger() {

	encoder := getEncoder()

	writeSyncer := getWriteSyncer()

	levelEnabler := getLevelEnabler()

	consoleEncoder := getConsoleEncoder()

	newCore := zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, levelEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
	)

	Logger = zap.New(newCore, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
}

// GetEncoder 自定义Encoder
func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,      // 默认换行符"\n"
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 日志等级序列为小写字符串，如:InfoLevel被序列化为 "info"
			EncodeTime:     zapcore.EpochTimeEncoder,       // 日志时间格式显示
			EncodeDuration: zapcore.SecondsDurationEncoder, // 时间序列化，Duration为经过的浮点秒数
			EncodeCaller:   zapcore.ShortCallerEncoder,     // 日志行号显示
		})
}

// getConsoleEncoder 输出日志到控制台
func getConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// GetWriteSyncer 自定义WriteSyncer
func getWriteSyncer() zapcore.WriteSyncer {
	lumberjack := &lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     30,
	}
	return zapcore.AddSync(lumberjack)
}

// GetLevelEnabler 自定义的LevelEnabler
func getLevelEnabler() zapcore.Level {
	return zapcore.InfoLevel // 只会打印出info及其以上级别的日志
}
