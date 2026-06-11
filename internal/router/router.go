package router

import (
	"errors"
	"net/http"
	"strconv"

	"GoStudy/internal/auth"
	"GoStudy/internal/middleware"
	"GoStudy/internal/response"
	"GoStudy/internal/service"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AppName      string
	Env          string
	TokenManager *auth.Manager
	UserService  *service.UserService
	TaskService  *service.TaskService
}

type Router struct {
	engine *gin.Engine
	deps   Dependencies
}

// New 创建 Gin 引擎并注册企业骨架的公共中间件和路由。
func New(deps Dependencies) *Router {
	r := &Router{
		engine: gin.New(),
		deps:   deps,
	}
	r.engine.Use(middleware.RequestID(), middleware.CORS(), gin.Recovery(), middleware.RequestLogger())
	r.register()
	return r
}

// Engine 暴露 http.Handler，方便应用层托管 Server 和测试层发起请求。
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// register 区分公开接口和受保护接口，避免鉴权逻辑散落在处理函数里。
func (r *Router) register() {
	r.engine.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{
			"name": r.deps.AppName,
			"env":  r.deps.Env,
		})
	})

	api := r.engine.Group("/api/v1")
	{
		api.POST("/users", r.createUser)
		api.POST("/auth/login", r.login)

		protected := api.Group("")
		protected.Use(middleware.Auth(r.deps.TokenManager))
		{
			protected.GET("/users", r.listUsers)
			protected.POST("/tasks", r.createTask)
			protected.GET("/tasks", r.listTasks)
			protected.PATCH("/tasks/:id/complete", r.completeTask)
		}
	}

	r.engine.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "route not found")
	})
}

// createUser 负责 HTTP 入参绑定，业务规则交给 service 处理。
func (r *Router) createUser(c *gin.Context) {
	var req service.RegisterUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := r.deps.UserService.Register(c.Request.Context(), req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Created(c, user)
}

// login 校验账号密码并签发接口访问 Token。
func (r *Router) login(c *gin.Context) {
	var req service.LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := r.deps.UserService.Login(c.Request.Context(), req)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	token, err := r.deps.TokenManager.Generate(user.ID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// listUsers 返回用户列表，当前属于受保护接口。
func (r *Router) listUsers(c *gin.Context) {
	users, err := r.deps.UserService.List(c.Request.Context())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, users)
}

// createTask 创建任务，任务归属校验放在服务层完成。
func (r *Router) createTask(c *gin.Context) {
	var req service.CreateTaskInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := r.deps.TaskService.Create(c.Request.Context(), req)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Created(c, task)
}

// listTasks 返回任务列表。
func (r *Router) listTasks(c *gin.Context) {
	tasks, err := r.deps.TaskService.List(c.Request.Context())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, tasks)
}

// completeTask 将路径参数转成业务 ID，再调用服务层完成任务。
func (r *Router) completeTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := r.deps.TaskService.Complete(c.Request.Context(), id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, task)
}

// writeServiceError 将业务错误稳定映射为 HTTP 状态码。
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrConflict):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
