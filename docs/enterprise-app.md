# 企业应用骨架

这个目录说明 `cmd/server` 示例的目标和运行方式。它不是替代原有学习示例，而是在仓库里新增一个更接近企业后端的组合示例。

## 当前能力

- `cmd/server`：应用启动入口。
- `internal/app`：装配配置、仓储、服务和路由。
- `internal/router`：Gin 路由、版本分组和错误响应。
- `internal/service`：用户和任务业务逻辑。
- `internal/repository`：仓储接口、内存实现和 GORM/MySQL 实现。
- `internal/domain`：用户、任务等领域对象。
- `internal/auth`：演示用密码摘要和 Token 签发校验。
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
| `internal/domain` | 保存业务对象，不依赖 Gin、GORM 等框架 |

这样的拆分让示例可以先用内存仓储学习流程，再平滑切换到 MySQL。

## 配置项

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `APP_NAME` | `GoStudy Enterprise` | 应用名称 |
| `APP_ENV` | `dev` | 运行环境，非 `prod` 时日志也会输出到控制台 |
| `APP_ADDR` | `:8080` | HTTP 监听地址 |
| `APP_STORAGE` | `memory` | 仓储类型，支持 `memory`、`mysql` |
| `APP_TOKEN_SECRET` | `gostudy-dev-secret` | Token 签名密钥 |
| `APP_TOKEN_TTL` | `2h` | Token 有效期 |
| `APP_SHUTDOWN_TIMEOUT` | `10s` | 优雅停机等待时间 |
| `MYSQL_DSN` | 空 | MySQL 连接串，`APP_STORAGE=mysql` 时必填 |

## 运行方式

```powershell
$env:APP_NAME="GoStudy Enterprise"
$env:APP_ENV="dev"
$env:APP_ADDR=":8080"
$env:APP_STORAGE="memory"
$env:APP_TOKEN_SECRET="change-me"
$env:APP_TOKEN_TTL="2h"
$env:APP_SHUTDOWN_TIMEOUT="10s"
go run ./cmd/server
```

切换到 MySQL：

```powershell
$env:APP_STORAGE="mysql"
$env:MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
go run ./cmd/server
```

MySQL 模式会使用 GORM 自动迁移用户表和任务表。没有启动 MySQL 时，保持默认 `APP_STORAGE=memory` 即可运行全部本地测试。

## 接口清单

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| `GET` | `/health` | 否 | 健康检查 |
| `POST` | `/api/v1/users` | 否 | 注册用户 |
| `POST` | `/api/v1/auth/login` | 否 | 登录并取得 Token |
| `GET` | `/api/v1/users` | 是 | 查询用户列表 |
| `POST` | `/api/v1/tasks` | 是 | 创建任务 |
| `GET` | `/api/v1/tasks` | 是 | 查询任务列表 |
| `PATCH` | `/api/v1/tasks/:id/complete` | 是 | 完成任务 |

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
  -d '{"title":"发布企业骨架","owner_id":1}'
```

完成任务：

```sh
curl -X PATCH http://127.0.0.1:8080/api/v1/tasks/1/complete \
  -H "Authorization: Bearer <token>"
```

## 后续升级方向

1. 接入 Redis 缓存用户或任务查询。
2. 接入 RabbitMQ 发布任务事件。
3. 将演示 Token 替换为 JWT 或统一认证中心。
4. 继续补充更多 HTTP 路由测试。

更完整的阶段拆分和验收方式见 `docs/enterprise-roadmap.md`。
