package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"GoStudy/internal/auth"
	"GoStudy/internal/cache"
	"GoStudy/internal/event"
	"GoStudy/internal/repository"
	"GoStudy/internal/router"
	"GoStudy/internal/service"
)

// Config 保存企业应用运行所需的基础配置。
type Config struct {
	Name            string
	Addr            string
	Env             string
	Storage         string
	MySQLDSN        string
	Cache           string
	RedisAddr       string
	RedisPassword   string
	MQ              string
	RabbitMQURL     string
	RabbitMQQueue   string
	TokenSecret     string
	TokenTTL        time.Duration
	ShutdownTimeout time.Duration
}

// App 组合配置、仓储、服务和路由，作为应用装配层。
type App struct {
	cfg    Config
	store  repository.Store
	cache  cache.Store
	events event.Publisher
	router *router.Router
}

// New 根据配置选择仓储实现，再装配服务层和路由层。
func New(cfg Config) (*App, error) {
	store, err := repository.NewStore(cfg.Storage, cfg.MySQLDSN)
	if err != nil {
		return nil, err
	}
	cacheStore, err := cache.New(cfg.Cache, cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		return nil, err
	}
	publisher, err := event.New(cfg.MQ, cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		_ = cacheStore.Close()
		return nil, err
	}

	tokenManager := auth.NewManager(cfg.TokenSecret, cfg.TokenTTL)

	userService := service.NewUserService(store, cacheStore)
	taskService := service.NewTaskService(store, store, publisher)

	return &App{
		cfg:    cfg,
		store:  store,
		cache:  cacheStore,
		events: publisher,
		router: router.New(router.Dependencies{
			AppName:      cfg.Name,
			Env:          cfg.Env,
			TokenManager: tokenManager,
			UserService:  userService,
			TaskService:  taskService,
			Health: func() map[string]string {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				return map[string]string{
					"storage": storageHealth(ctx, store, cfg.Storage),
					"cache":   cacheStore.Health(ctx),
					"mq":      publisher.Health(ctx),
				}
			},
		}),
	}, nil
}

// Run 托管 HTTP Server 生命周期，收到外部取消信号后执行优雅停机。
func (a *App) Run(ctx context.Context) error {
	defer a.closeStore()

	timeout := a.cfg.ShutdownTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	server := &http.Server{
		Addr:              a.cfg.Addr,
		Handler:           a.router.Engine(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

func (a *App) closeStore() {
	_ = a.cache.Close()
	_ = a.events.Close()
	closer, ok := a.store.(interface {
		Close() error
	})
	if ok {
		_ = closer.Close()
	}
}

func storageHealth(ctx context.Context, store repository.Store, storage string) string {
	if storage == "" || storage == repository.StorageMemory {
		return "skipped"
	}
	if checker, ok := store.(interface{ Health(context.Context) string }); ok {
		return checker.Health(ctx)
	}
	return "ok"
}
