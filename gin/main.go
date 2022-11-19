package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var R *gin.Engine

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`            // 用户名 必填
	Password string `form:"password" json:"password" xml:"password" binding:"required"` // 密码 必填
}

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func init() {
	//记录日志
	file, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r := gin.New()
	R = r

	// 跨域设置
	R.Use(Cors())

	//使用logger中间件
	//R.Use(ShowCode())
	R.Use(Logger())

	// CommonRouterRegister
	CommonRouterRegister()

	// FileUpload
	FileUpload()

	// RouterGroup
	RouterGroup()

	// Goroutine
	R.GET("/longAsync", longAsync)

	// Mode Bind
	R.POST("/modelBind", ModelBind)

	// Url Bind
	r.GET("/urlBind/:name/:id", UrlBindFunc)

	//Redirect
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})

	R.Run(":8080")
}

func main() {
	fmt.Println("welcome to gin!!")

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

// Logger 自定义中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 请求前

		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Print(latency)

		// 获取发送的 status
		status := c.Writer.Status()
		log.Println(status)
	}
}

func ShowCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("show me code")
	}
}

// UrlBindFunc url绑定
func UrlBindFunc(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
}

// ModelBind 模型绑定验证
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

// longAsync Goroutine使用
func longAsync(c *gin.Context) {
	cCp := c.Copy()
	go func() {
		time.Sleep(5 * time.Second)
		log.Println(" request path = " + cCp.Request.URL.Path)
	}()
}

// login
func login(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

// submit
func submit(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

// RouterGroup 认证路由组
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

// CommonRouterRegister 普通路由注册
func CommonRouterRegister() {
	// gin.H 是 map[string]interface{} 的一种快捷方式
	R.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	R.GET("/moreJSON", func(c *gin.Context) {
		// 你也可以使用一个结构体
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// 注意 msg.Name 在 JSON 中变成了 "user"
		// 将输出：{"user": "Lena", "Message": "hey", "Number": 123}
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
		// protobuf 的具体定义写在 testdata/protoExample 文件中。
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		// 请注意，数据在响应中变为二进制数据
		// 将输出被 protoExample.Test protobuf 序列化了的数据
		c.ProtoBuf(http.StatusOK, data)
	})
}

// FileUpload 文件上传路由
func FileUpload() {
	// 单文件上传
	R.MaxMultipartMemory = 8 << 20 // 8 MiB
	R.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "upload/" + file.Filename
		// 上传文件至指定的完整文件路径
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// 多文件上传
	R.MaxMultipartMemory = 8 << 20 // 8 MiB
	R.POST("/uploads", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			dst := "upload/" + file.Filename
			// 上传文件至指定目录
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
}
