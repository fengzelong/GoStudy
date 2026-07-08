package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"GoStudy/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
)

var R *gin.Engine

// ginLogFile 由 setupGinLog 持有，供 main 退出前 Close。
var ginLogFile *os.File

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`             // 用户名，必填
	Password string `form:"password" json:"password" xml:"password" binding:"required"` // 密码，必填
}

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

type QueryInput struct {
	Keyword  string `form:"keyword" binding:"required"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type FormInput struct {
	Name    string `form:"name" binding:"required"`
	Message string `form:"message" binding:"required"`
}

type HeaderInput struct {
	RequestID string `header:"X-Request-Id" binding:"required"`
}

func setupRouter() *gin.Engine {
	setupGinLog()

	r := gin.New()
	r.HandleMethodNotAllowed = true
	R = r

	R.Use(gin.Recovery())
	R.Use(Cors())
	R.Use(Logger())

	R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "路由不存在"})
	})
	R.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "请求方法不允许"})
	})

	RegisterBasicRoutes(R)
	RegisterBindRoutes(R)
	RegisterResponseRoutes(R)
	RegisterFileRoutes(R)
	RegisterRouterGroup(R)
	RegisterAuthRoutes(R)

	return R
}

func main() {
	addr := config.Env("GIN_ADDR", ":8080")
	fmt.Println("welcome to gin!!")

	srv := &http.Server{
		Addr:              addr,
		Handler:           setupRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	defer func() {
		if ginLogFile != nil {
			_ = ginLogFile.Close()
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gin 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Gin 服务关闭失败: %v", err)
	}
}

func setupGinLog() {
	logFile := config.Env("GIN_LOG_FILE", "gin.log")
	if logFile == "-" {
		gin.DefaultWriter = os.Stdout
		return
	}

	// 用追加模式打开，避免每次启动截断历史日志。
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("创建 Gin 日志文件失败: %v", err)
		gin.DefaultWriter = os.Stdout
		return
	}
	ginLogFile = file

	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 凭证请求需要回显具体 Origin；Access-Control-Allow-Origin: * 与 Allow-Credentials: true 同时出现会被浏览器拒绝。
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,X-Request-Id")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Content-Type", "application/json")
		}
		if method == http.MethodOptions && origin != "" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Logger 自定义中间件示例：把请求日志写到 gin.DefaultWriter，
// 跟随 setupGinLog 的配置（文件 + stdout 或纯 stdout）。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "12345")
		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()
		fmt.Fprintf(gin.DefaultWriter, "path=%s status=%d latency=%s\n", c.Request.URL.Path, status, latency)
	}
}

func ShowCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("show_code", "show me code")
		c.Next()
	}
}

// RegisterBasicRoutes 注册基础路由、查询参数、Header 和 Cookie 示例。
func RegisterBasicRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/param/:name", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"name": c.Param("name")})
	})

	r.GET("/query", func(c *gin.Context) {
		name := c.DefaultQuery("name", "jack")
		age := c.Query("age")
		c.JSON(http.StatusOK, gin.H{"name": name, "age": age})
	})

	r.GET("/header", func(c *gin.Context) {
		c.Header("X-Gin-Example", "header")
		c.JSON(http.StatusOK, gin.H{"user_agent": c.GetHeader("User-Agent")})
	})

	r.GET("/cookie/set", func(c *gin.Context) {
		c.SetCookie("gin_user", "jack", 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "cookie set"})
	})

	r.GET("/cookie/read", func(c *gin.Context) {
		value, err := c.Cookie("gin_user")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cookie 不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"gin_user": value})
	})

	r.GET("/middleware/code", ShowCode(), func(c *gin.Context) {
		code, _ := c.Get("show_code")
		c.JSON(http.StatusOK, gin.H{"code": code})
	})

	r.GET("/longAsync", longAsync)
	r.GET("/redirect", func(c *gin.Context) {
		// 用 307 而非 301：301 会被浏览器永久缓存，无法再改目标。
		c.Redirect(http.StatusTemporaryRedirect, "https://gin-gonic.com/")
	})
}

// RegisterBindRoutes 注册参数绑定和校验示例。
func RegisterBindRoutes(r *gin.Engine) {
	r.GET("/urlBind/:name/:id", UrlBindFunc)
	r.POST("/modelBind", ModelBind)

	r.GET("/bindQuery", func(c *gin.Context) {
		var input QueryInput
		if err := c.ShouldBindQuery(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Page == 0 {
			input.Page = 1
		}
		if input.PageSize == 0 {
			input.PageSize = 10
		}
		c.JSON(http.StatusOK, gin.H{
			"keyword":   input.Keyword,
			"page":      input.Page,
			"page_size": input.PageSize,
		})
	})

	r.POST("/form", func(c *gin.Context) {
		var input FormInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": input.Name, "message": input.Message})
	})

	r.GET("/bindHeader", func(c *gin.Context) {
		var input HeaderInput
		if err := c.ShouldBindHeader(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"request_id": input.RequestID})
	})

	r.POST("/bindXML", func(c *gin.Context) {
		var input Login
		if err := c.ShouldBindXML(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": input.User})
	})
}

// RegisterResponseRoutes 注册常见响应示例。
func RegisterResponseRoutes(r *gin.Engine) {
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		c.ProtoBuf(http.StatusOK, data)
	})

	r.GET("/download", func(c *gin.Context) {
		c.Header("Content-Disposition", "attachment; filename=gin.txt")
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("Gin 文件下载示例\n"))
	})

	r.GET("/stream", func(c *gin.Context) {
		step := 0
		c.Stream(func(w io.Writer) bool {
			step++
			if step > 3 {
				return false
			}
			c.SSEvent("message", gin.H{"step": step})
			return true
		})
	})
}

// RegisterFileRoutes 注册文件上传和静态文件路由。
func RegisterFileRoutes(r *gin.Engine) {
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.StaticFS("/static", http.Dir(uploadDir()))

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请上传 file 字段"})
			return
		}

		filename := filepath.Base(file.Filename)
		dst := filepath.Join(uploadDir(), filename)
		if err := os.MkdirAll(uploadDir(), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", filename))
	})

	r.POST("/uploads", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请上传 upload[] 字段"})
			return
		}

		files := form.File["upload[]"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请上传 upload[] 字段"})
			return
		}

		if err := os.MkdirAll(uploadDir(), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, file := range files {
			filename := filepath.Base(file.Filename)
			dst := filepath.Join(uploadDir(), filename)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
}

// RegisterRouterGroup 注册分组路由。
func RegisterRouterGroup(r *gin.Engine) {
	authorizedV1 := r.Group("/v1")
	{
		authorizedV1.POST("/login", login)
		authorizedV1.POST("/submit", submit)
	}
	authorizedV2 := r.Group("/v2")
	{
		authorizedV2.POST("/login", login)
	}
}

// RegisterAuthRoutes 注册 BasicAuth 示例。
func RegisterAuthRoutes(r *gin.Engine) {
	authorized := r.Group("/basic", gin.BasicAuth(gin.Accounts{
		"gin": "123456",
	}))
	authorized.GET("/secret", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "Gin BasicAuth"})
	})
}

// UrlBindFunc 绑定 URL 参数。
func UrlBindFunc(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": person.Name, "uuid": person.ID})
}

// ModelBind 校验 JSON 请求数据。
func ModelBind(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if json.User != "jack" || json.Password != "123" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

// longAsync 演示在处理函数中使用 goroutine。
func longAsync(c *gin.Context) {
	cCp := c.Copy()
	go func() {
		time.Sleep(5 * time.Second)
		log.Println("request path = " + cCp.Request.URL.Path)
	}()
	c.JSON(http.StatusAccepted, gin.H{"status": "async task started"})
}

func login(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(http.StatusOK, fmt.Sprintf("hello %s\n", name))
}

func submit(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(http.StatusOK, fmt.Sprintf("hello %s\n", name))
}

func uploadDir() string {
	return config.Env("GIN_UPLOAD_DIR", filepath.Join("gin", "upload"))
}
