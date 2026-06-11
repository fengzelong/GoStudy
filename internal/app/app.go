package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"GoStudy/internal/auth"
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
	TokenSecret     string
	TokenTTL        time.Duration
	ShutdownTimeout time.Duration
}

// App 组合配置、仓储、服务和路由，作为应用装配层。
type App struct {
	cfg    Config
	store  repository.Store
	router *router.Router
}

// New 根据配置选择仓储实现，再装配服务层和路由层。
func New(cfg Config) (*App, error) {
	store, err := repository.NewStore(cfg.Storage, cfg.MySQLDSN)
	if err != nil {
		return nil, err
	}

	tokenManager := auth.NewManager(cfg.TokenSecret, cfg.TokenTTL)

	userService := service.NewUserService(store)
	taskService := service.NewTaskService(store, store)

	return &App{
		cfg:   cfg,
		store: store,
		router: router.New(router.Dependencies{
			AppName:      cfg.Name,
			Env:          cfg.Env,
			TokenManager: tokenManager,
			UserService:  userService,
			TaskService:  taskService,
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
	closer, ok := a.store.(interface {
		Close() error
	})
	if ok {
		_ = closer.Close()
	}
}
