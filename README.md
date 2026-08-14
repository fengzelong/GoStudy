# GoStudy

GoStudy 是一个 Go 语言学习示例仓库，覆盖基础语法、算法、并发、网络编程、Web 框架、中间件、设计模式和一个可运行的分层企业应用骨架。

## 目录说明

| 主题 | 目录 | 内容 |
| --- | --- | --- |
| 算法 | [`algorithm`](algorithm) | [`binsearch`](algorithm/binsearch)、[`bubblesort`](algorithm/bubblesort)、[`heapsort`](algorithm/heapsort)、[`quicksort`](algorithm/quicksort)、[`selectsort`](algorithm/selectsort) |
| 基础语法 | [`array`](array) | 数组、切片、映射、结构体、函数和递归 |
| 设计模式 | [`designpattern`](designpattern) | 13 种常见设计模式的 Go 实现 |
| 并发 | [`goroutine`](goroutine)、[`sync`](sync) | goroutine、channel、工作池、原子操作、`sync.Map`、读写锁 |
| Web | [`gin`](gin)、[`http`](http) | Gin 路由和中间件；标准库 HTTP 示例 |
| 网络 | [`tcp`](tcp)、[`udp`](udp)、[`rpc`](rpc)、[`websocket`](websocket) | 回显、标准库 RPC 与 WebSocket 示例 |
| gRPC | [`grpc`](grpc)、[`proto`](proto) | 客户端、服务端、protobuf 协议与生成代码 |
| 数据与消息 | [`gomysql`](gomysql)、[`gorm`](gorm)、[`goredis`](goredis)、[`gorabbitmq`](gorabbitmq) | MySQL、GORM、Redis、RabbitMQ；需要对应本地服务 |
| 其他专题 | [`encrypt`](encrypt)、[`log`](log)、[`runtime`](runtime)、[`task`](task)、[`gmp`](gmp)、[`study`](study) | 加密、日志、运行时、定时任务、trace 与语法练习 |
| 企业骨架 | [`cmd/server`](cmd/server)、[`internal`](internal) | 分层、鉴权、分页、任务归属、可选缓存与事件发布 |

## 模块说明

仓库根目录是名为 `GoStudy` 的 Go 模块。以下示例保留独立的 `go.mod`，方便单独运行：

- `gomysql`
- `gorabbitmq`
- `goredis`
- `gorm`
- `websocket`

## 常用命令

运行完整回归测试：

```powershell
.\scripts\test.ps1
```

提交前快速检查：

```powershell
.\scripts\check.ps1
```

运行根模块测试：

```sh
go test ./...
```

运行企业应用骨架：

```sh
go run ./cmd/server
```

启动后可访问：

```sh
curl http://127.0.0.1:8080/health
```

企业应用骨架默认使用内存仓储，缓存和消息队列均关闭，不依赖外部服务；需要连接 MySQL 时设置：

```powershell
$env:APP_STORAGE="mysql"
$env:MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
go run ./cmd/server
```

可选启用进程内缓存和任务事件，仍不需要外部服务：

```powershell
$env:APP_CACHE="memory"
$env:APP_MQ="memory"
go run ./cmd/server
```

完整的 Redis、RabbitMQ 配置和接口说明见 [企业应用骨架说明](docs/enterprise-app.md)。

运行单个示例：

```sh
go run ./gin
go run ./tcp/server
go run ./tcp/client
```

常见网络示例可以通过环境变量调整监听地址或连接地址，例如：

```powershell
$env:GIN_ADDR=":8081"
$env:TCP_ADDR="127.0.0.1:20001"
$env:GRPC_ADDR="127.0.0.1:50052"
```

如果 PowerShell 执行策略阻止运行脚本，可使用根模块测试命令兜底：

```powershell
$env:GOMODCACHE="$PWD\.tmp_gomodcache"
$env:GOPATH="$PWD\.tmp_gopath"
go test ./...
```

运行独立子模块示例：

```sh
cd websocket
go run .
```

## 文档与规划

- [项目迭代计划](docs/iteration-plan.md)：全仓库的文档、测试和示例整理计划。
- [企业应用骨架说明](docs/enterprise-app.md)：企业服务的分层、配置、接口和运行方式。
- [OpenAPI 文档](docs/openapi.yaml)：企业服务接口、鉴权、分页和错误响应约定。
- [企业应用迭代规划](docs/enterprise-roadmap.md)：企业服务后续功能规划。
- [变更记录](CHANGELOG.md)：代表性迭代记录。

## 注意事项

- 外部服务连接示例见根目录 `.env.example`。
- 使用 `docker compose up -d` 可启动 MySQL、Redis、RabbitMQ 本地依赖；具体配置见企业应用骨架说明。
- Git 会忽略生成的二进制文件、日志、trace 文件、本地缓存、IDE 配置和运行时上传文件。
- 部分示例需要先启动 MySQL、Redis、RabbitMQ 或配套的网络客户端/服务端。
- 协作与编码约定见 `AGENTS.md`。
