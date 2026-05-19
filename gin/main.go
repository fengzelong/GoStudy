package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"GoStudy/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
)

var R *gin.Engine

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`            // 用户名，必填
	Password string `form:"password" json:"password" xml:"password" binding:"required"` // 密码，必填
}

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func setupRouter() *gin.Engine {
	file, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r := gin.New()
	R = r

	R.Use(Cors())
	R.Use(Logger())

	CommonRouterRegister()
	FileUpload()
	RouterGroup()

	R.GET("/longAsync", longAsync)
	R.POST("/modelBind", ModelBind)
	r.GET("/urlBind/:name/:id", UrlBindFunc)
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})

	return R
}

func main() {
	addr := config.Env("GIN_ADDR", ":8080")
	fmt.Println("welcome to gin!!")
	setupRouter().Run(addr)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// Logger 自定义中间件示例。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "12345")
		c.Next()

		latency := time.Since(t)
		log.Print(latency)

		status := c.Writer.Status()
		log.Println(status)
	}
}

func ShowCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("show me code")
	}
}

// UrlBindFunc 绑定 URL 参数。
func UrlBindFunc(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
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
		log.Println(" request path = " + cCp.Request.URL.Path)
	}()
}

func login(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

func submit(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

// RouterGroup 注册分组路由。
func RouterGroup() {
	authorizedV1 := R.Group("/v1")
	{
		authorizedV1.POST("/login", login)
		authorizedV1.POST("/submit", submit)
	}
	authorizedV2 := R.Group("/v2")
	{
		authorizedV2.POST("/login", login)
	}
}

// CommonRouterRegister 注册常见响应示例。
func CommonRouterRegister() {
	R.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	R.GET("/moreJSON", func(c *gin.Context) {
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

	R.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	R.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	R.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		c.ProtoBuf(http.StatusOK, data)
	})
}

// FileUpload 注册文件上传路由。
func FileUpload() {
	R.MaxMultipartMemory = 8 << 20 // 8 MiB
	R.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "upload/" + file.Filename
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	R.MaxMultipartMemory = 8 << 20 // 8 MiB
	R.POST("/uploads", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			dst := "upload/" + file.Filename
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
}
