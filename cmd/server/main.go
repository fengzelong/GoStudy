package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoStudy/internal/app"
	"GoStudy/internal/config"
	"GoStudy/internal/logger"
)

func main() {
	// 启动入口只负责读取配置、初始化基础设施和托管应用生命周期。
	tokenTTL, err := time.ParseDuration(config.Env("APP_TOKEN_TTL", "2h"))
	if err != nil {
		log.Fatalf("parse APP_TOKEN_TTL: %v", err)
	}
	shutdownTimeout, err := time.ParseDuration(config.Env("APP_SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		log.Fatalf("parse APP_SHUTDOWN_TIMEOUT: %v", err)
	}

	cfg := app.Config{
		Name:            config.Env("APP_NAME", "GoStudy Enterprise"),
		Addr:            config.Env("APP_ADDR", ":8080"),
		Env:             config.Env("APP_ENV", "dev"),
		Storage:         config.Env("APP_STORAGE", "memory"),
		MySQLDSN:        config.Env("MYSQL_DSN", ""),
		TokenSecret:     config.Env("APP_TOKEN_SECRET", "gostudy-dev-secret"),
		TokenTTL:        tokenTTL,
		ShutdownTimeout: shutdownTimeout,
	}

	closer, err := logger.Init(cfg.Env) 
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer closer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	if err := server.Run(ctx); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
