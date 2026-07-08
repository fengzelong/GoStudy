# gin

这个目录演示 Gin 常用能力，覆盖路由、参数绑定、响应格式、中间件、文件上传、Cookie、BasicAuth、流式响应和优雅关闭。

## 运行方式

```powershell
$env:GIN_ADDR=":8080"
$env:GIN_UPLOAD_DIR="gin/upload"
go run ./gin
```

如果不设置环境变量，服务默认监听 `:8080`，上传文件默认保存到 `gin/upload`。

## 路由示例

基础路由：

```sh
curl http://127.0.0.1:8080/ping
curl http://127.0.0.1:8080/param/tom
curl "http://127.0.0.1:8080/query?name=tom&age=18"
```

参数绑定：

```sh
curl -X POST http://127.0.0.1:8080/modelBind \
  -H "Content-Type: application/json" \
  -d '{"user":"jack","password":"123"}'

curl http://127.0.0.1:8080/urlBind/tom/550e8400-e29b-41d4-a716-446655440000
curl "http://127.0.0.1:8080/bindQuery?keyword=gin&page=1&page_size=10"
curl http://127.0.0.1:8080/bindHeader -H "X-Request-Id: req-1"
```

表单和 XML：

```sh
curl -X POST http://127.0.0.1:8080/form \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=tom&message=hello"

curl -X POST http://127.0.0.1:8080/bindXML \
  -H "Content-Type: application/xml" \
  -d '<Login><user>jack</user><password>123</password></Login>'
```

响应格式：

```sh
curl http://127.0.0.1:8080/someJSON
curl http://127.0.0.1:8080/someXML
curl http://127.0.0.1:8080/someYAML
curl http://127.0.0.1:8080/someProtoBuf
curl http://127.0.0.1:8080/download
curl http://127.0.0.1:8080/stream
```

Cookie、Header 和 BasicAuth：

```sh
curl -i http://127.0.0.1:8080/header
curl -i http://127.0.0.1:8080/cookie/set
curl http://127.0.0.1:8080/cookie/read --cookie "gin_user=jack"
curl -u gin:123456 http://127.0.0.1:8080/basic/secret
```

文件上传：

```sh
curl -X POST http://127.0.0.1:8080/upload \
  -F "file=@README.md"

curl -X POST http://127.0.0.1:8080/uploads \
  -F "upload[]=@README.md" \
  -F "upload[]=@go.mod"
```

路由分组：

```sh
curl -X POST "http://127.0.0.1:8080/v1/login?name=jack"
curl -X POST "http://127.0.0.1:8080/v1/submit?name=lily"
curl -X POST "http://127.0.0.1:8080/v2/login?name=tom"
```

## 验证

```powershell
go test ./gin
```
