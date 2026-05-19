# GoStudy

GoStudy 是一个 Go 语言学习示例仓库，覆盖基础语法、算法、并发、网络编程、
Web 框架、中间件和常见设计模式。

## 目录说明

| 目录 | 内容 |
| --- | --- |
| `algorithm` | 二分查找和常见排序示例 |
| `array` | 基础语法、切片、映射、结构体和基准测试 |
| `designpattern` | 常见设计模式的 Go 实现 |
| `gin` | Gin 路由、参数绑定、中间件、文件上传和响应示例 |
| `goroutine` | goroutine、channel、select、WaitGroup、协程池、map 和 atomic 示例 |
| `grpc`, `proto` | gRPC 客户端、服务端和 protobuf 生成代码 |
| `http`, `tcp`, `udp`, `rpc`, `websocket` | 网络编程示例 |
| `gomysql`, `goredis`, `gorm`, `gorabbitmq` | 数据库和中间件示例 |
| `log`, `encrypt`, `runtime`, `sync`, `task`, `gmp`, `study` | 其他专题示例 |

## 模块说明

仓库根目录是名为 `GoStudy` 的 Go 模块。部分中间件示例保留了独立的
`go.mod`，便于单独运行：

- `gomysql`
- `gorabbitmq`
- `goredis`
- `gorm`
- `websocket`

## 常用命令

运行根模块测试：

```sh
go test ./...
```

运行单个示例：

```sh
go run ./gin
go run ./tcp/server
go run ./tcp/client
```

运行独立子模块示例：

```sh
cd websocket
go run .
```

## 注意事项

- Git 会忽略生成的二进制文件、日志、trace 文件、本地缓存、IDE 配置和运行时上传文件。
- 部分示例需要先启动本地服务，例如 MySQL、Redis、RabbitMQ，或配套的网络客户端/服务端。
- 协作和编码氛围约定见 `AGENTS.md`。
