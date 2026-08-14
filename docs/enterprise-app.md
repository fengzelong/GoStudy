# 企业应用骨架

这个目录说明 `cmd/server` 示例的目标和运行方式。它不是替代原有学习示例，而是在仓库里新增一个更接近企业后端的组合示例。

## 当前能力

- `cmd/server`：应用启动入口。
- `internal/app`：装配配置、仓储、服务和路由。
- `internal/router`：Gin 路由、版本分组、分页参数和错误响应。
- `internal/service`：用户和任务业务逻辑、任务归属隔离和分页。
- `internal/repository`：仓储接口、内存实现和 GORM/MySQL 实现。
- `internal/cache`：关闭、内存和 Redis 缓存实现，当前缓存用户资料。
- `internal/event`：关闭、内存和 RabbitMQ 事件发布器，当前发布任务创建、完成事件。
- `internal/domain`：用户、任务等领域对象。
- `internal/auth`：bcrypt 密码摘要和 HS256 JWT 签发校验。
- `internal/audit`：登录、任务创建和完成的内存审计记录。
- `internal/logger`：zap 结构化日志和日志切割。
- `internal/middleware`：请求 ID、CORS、请求日志和鉴权中间件。
- `internal/response`：统一 JSON 响应。
- 应用入口监听系统停止信号，收到中断后会优雅关闭 HTTP Server。

## 分层说明

企业骨架采用从外到内的分层方式：

| 层级 | 职责 |
| --- | --- |
| `cmd/server` | 读取环境变量，初始化日志，监听退出信号 |
| `internal/app` | 组合仓储、服务、路由和 HTTP Server 生命周期 |
| `internal/router` | 处理 HTTP 入参、鉴权分组和响应转换 |
| `internal/service` | 承载业务规则，例如注册去重、任务归属校验 |
| `internal/repository` | 隐藏存储细节，支持内存和 MySQL 两种实现 |
| `internal/cache` | 隔离缓存实现，支持关闭、内存和 Redis 三种模式 |
| `internal/event` | 隔离事件发布实现，支持关闭、内存和 RabbitMQ 三种模式 |
| `internal/domain` | 保存业务对象，不依赖 Gin、GORM 等框架 |
| `internal/auth` | 使用 bcrypt 校验密码并签发包含角色声明的 HS256 JWT |
| `internal/audit` | 记录登录与任务状态变更等关键操作 |

这样的拆分让示例可以先用内存仓储学习流程，再平滑切换到 MySQL。

## 配置项

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `APP_NAME` | `GoStudy Enterprise` | 应用名称 |
| `APP_ENV` | `dev` | 运行环境，非 `prod` 时日志也会输出到控制台 |
| `APP_ADDR` | `:8080` | HTTP 监听地址 |
| `APP_STORAGE` | `memory` | 仓储类型，支持 `memory`、`mysql` |
| `APP_CACHE` | `off` | 缓存类型，支持 `off`、`memory`、`redis` |
| `APP_MQ` | `off` | 事件发布类型，支持 `off`、`memory`、`rabbitmq` |
| `APP_TOKEN_SECRET` | `gostudy-dev-secret` | Token 签名密钥 |
| `APP_TOKEN_TTL` | `2h` | Token 有效期 |
| `APP_ADMIN_EMAIL` | 空 | 匹配该邮箱的新注册用户或既有用户会被设为管理员 |
| `APP_SHUTDOWN_TIMEOUT` | `10s` | 优雅停机等待时间 |
| `MYSQL_DSN` | 空 | MySQL 连接串，`APP_STORAGE=mysql` 时必填 |
| `REDIS_ADDR` | 空 | Redis 地址，`APP_CACHE=redis` 时必填 |
| `REDIS_PASSWORD` | 空 | Redis 密码 |
| `RABBITMQ_URL` | 空 | RabbitMQ 地址，`APP_MQ=rabbitmq` 时必填 |
| `RABBITMQ_QUEUE` | 空 | RabbitMQ 队列名，`APP_MQ=rabbitmq` 时必填 |

## 运行方式

```powershell
$env:APP_NAME="GoStudy Enterprise"
$env:APP_ENV="dev"
$env:APP_ADDR=":8080"
$env:APP_STORAGE="memory"
$env:APP_CACHE="off"
$env:APP_MQ="off"
$env:APP_TOKEN_SECRET="change-me"
$env:APP_TOKEN_TTL="2h"
$env:APP_ADMIN_EMAIL="admin@example.com"
$env:APP_SHUTDOWN_TIMEOUT="10s"
go run ./cmd/server
```

## 可选 MySQL 集成测试

企业骨架的 GORM 仓储集成测试默认跳过。它会创建、更新、查询并清理测试用户和任务；仅应连接专用测试库，例如 Compose 默认创建的 `go_test`。

先启动 MySQL 依赖，再显式开启测试：

```powershell
docker compose up -d mysql
$env:APP_INTEGRATION_MYSQL="1"
$env:MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
go test ./internal/repository -run TestGormStoreIntegration -count=1
```

不设置 `APP_INTEGRATION_MYSQL=1` 时，默认的 `go test ./...` 和 `scripts/test.ps1` 不会连接 MySQL。

切换到 MySQL：

```powershell
$env:APP_STORAGE="mysql"
$env:MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
go run ./cmd/server
```

MySQL 模式会使用 GORM 自动迁移用户表和任务表。没有启动 MySQL 时，保持默认 `APP_STORAGE=memory` 即可运行全部本地测试。

缓存和消息队列默认关闭，因此也不依赖 Redis 或 RabbitMQ。需要在本地观察缓存和事件行为时，可以使用不需要外部服务的内存实现：

```powershell
$env:APP_CACHE="memory"
$env:APP_MQ="memory"
go run ./cmd/server
```

连接外部 Redis 和 RabbitMQ 时显式启用对应模式：

```powershell
$env:APP_CACHE="redis"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:APP_MQ="rabbitmq"
$env:RABBITMQ_URL="amqp://guest:guest@127.0.0.1:5672/"
$env:RABBITMQ_QUEUE="gostudy.task.events"
go run ./cmd/server
```

`GET /health` 的 `data.dependencies` 会返回 `storage`、`cache` 和 `mq` 的状态。关闭或内存模式没有外部依赖时会返回 `skipped` 或 `ok`。

## Docker Compose 本地依赖

仓库根目录的 [`docker-compose.yml`](../docker-compose.yml) 提供 MySQL 8、Redis 7 和
RabbitMQ 3（含管理界面）编排。它只启动依赖服务，不会启动 Go 应用，也不会改变默认的
内存运行模式。

启动依赖：

```sh
docker compose up -d
```

查看状态和日志：

```sh
docker compose ps
docker compose logs -f mysql redis rabbitmq
```

停止并保留数据卷：

```sh
docker compose down
```

如需同时删除本地容器数据卷：

```sh
docker compose down -v
```

容器启动后，RabbitMQ 管理界面为 `http://127.0.0.1:15672`，默认账号密码均为 `guest`。
启用完整外部依赖的应用配置：

```powershell
$env:APP_STORAGE="mysql"
$env:MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
$env:APP_CACHE="redis"
$env:REDIS_ADDR="127.0.0.1:6379"
$env:APP_MQ="rabbitmq"
$env:RABBITMQ_URL="amqp://guest:guest@127.0.0.1:5672/"
$env:RABBITMQ_QUEUE="gostudy.task.events"
go run ./cmd/server
```

## 接口清单

机器可读的完整接口契约见 [OpenAPI 文档](openapi.yaml)。

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| `GET` | `/health` | 否 | 健康检查 |
| `POST` | `/api/v1/users` | 否 | 注册用户 |
| `POST` | `/api/v1/auth/login` | 否 | 登录并取得 Token |
| `GET` | `/api/v1/me` | 是 | 查询当前用户资料 |
| `GET` | `/api/v1/users` | 管理员 | 查询用户列表，支持分页 |
| `POST` | `/api/v1/tasks` | 是 | 创建当前用户的任务 |
| `GET` | `/api/v1/tasks` | 是 | 查询当前用户的任务，支持分页 |
| `PATCH` | `/api/v1/tasks/:id/complete` | 是 | 完成任务 |

列表接口支持 `page` 和 `page_size` 参数，默认分别为 `1` 和 `20`，`page_size` 最大为
`100`。响应的 `data` 包含 `items`、`page`、`page_size` 和 `total`。任务创建时不再接受
`owner_id`，任务自动归属当前认证用户；读取或完成其他用户的任务会返回 `403`。

普通注册用户的角色为 `user`。设置 `APP_ADMIN_EMAIL` 后，与该邮箱匹配的用户角色为
`admin`；当前管理员可以访问用户列表，普通用户会收到 `40301`。登录签发的是 HS256 JWT，
其中包含用户 ID、角色、签发时间和过期时间。密码以 bcrypt 摘要保存，接口不会返回密码摘要。

当前会将登录、任务创建和任务完成写入内存审计记录。审计接口和持久化存储将在后续工程化阶段扩展。

## 错误码

HTTP 状态码表示传输层结果，响应体 `code` 表示稳定的业务错误码：

| 业务码 | HTTP 状态 | 说明 |
| --- | --- | --- |
| `0` | `200`、`201` | 成功 |
| `40001` | `400` | 请求参数不合法 |
| `40101` | `401` | 未认证或 Token 无效 |
| `40301` | `403` | 无权访问当前资源 |
| `40401` | `404` | 资源或路由不存在 |
| `40901` | `409` | 资源冲突，例如邮箱已存在 |
| `50001` | `500` | 未预期的服务端错误 |

健康检查：

```sh
curl http://127.0.0.1:8080/health
```

创建用户：

```sh
curl -X POST http://127.0.0.1:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret1"}'
```

登录并取得 Token：

```sh
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret1"}'
```

创建任务：

```sh
curl -X POST http://127.0.0.1:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"发布企业骨架"}'
```

查询当前用户的第二页任务：

```sh
curl "http://127.0.0.1:8080/api/v1/tasks?page=2&page_size=20" \
  -H "Authorization: Bearer <token>"
```

完成任务：

```sh
curl -X PATCH http://127.0.0.1:8080/api/v1/tasks/1/complete \
  -H "Authorization: Bearer <token>"
```

## 后续升级方向

1. 为审计记录增加查询接口和持久化实现。
2. 根据多服务场景评估 JWT 刷新、吊销或统一认证中心。
3. 为 Redis、RabbitMQ 增加可选端到端集成测试。

更完整的阶段拆分和验收方式见 `docs/enterprise-roadmap.md`。
