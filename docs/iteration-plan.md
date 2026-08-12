# 项目迭代计划

更新时间：2026-08-12

这份规划覆盖 GoStudy 仓库里尚未整理或尚未补齐的部分，重点关注基础学习示例、
目录文档、测试覆盖和工程化检查。它与
[enterprise-roadmap.md](enterprise-roadmap.md) 互补：本文侧重全局清理与基础
示例，roadmap 文档继续负责企业骨架的纵向推进。

## 当前基线

### 已完成事项

- 已修正设计模式目录中的观察者模式拼写：`designpattern/observer`。
- 已补齐 `gin/` 常用示例，覆盖路由、参数绑定、响应格式、中间件、文件上传、
  Cookie、BasicAuth、SSE、404/405 和优雅关闭。
- 已新增 `gin/README.md`，列出环境变量、路由清单、curl 示例和测试命令。
- 已新增 `gin/main_test.go`，覆盖 ping、绑定、表单、Header、Cookie、Auth 和上传。
- 已在 `.env.example` 中补充 `GIN_LOG_FILE`、`GIN_UPLOAD_DIR`。
- 已完成第一阶段基线修正：统一 gRPC 端口示例、补充 PowerShell 执行策略兜底命令、
  `scripts/check.ps1` 增加 `go vet`、补充 `gmp/README.md`、清理
  `cmd/server/main.go` 尾随空白。
- 已补充第二阶段一批基础测试：`encrypt`、`log`、`runtime`、`task`、`http`、
  `goroutine`、`sync/map`、`sync/mutex`、`tcp`、`udp`、`grpc/server`、
  `grpc/client`、`rpc/client`、`rpc/server`、`study`。
- 已为独立子模块补充测试开关：`gomysql`、`goredis`、`gorm`、`gorabbitmq`
  在缺少外部服务环境变量时默认跳过真实连接测试，`websocket` 覆盖环境变量和
  用户列表处理函数。
- 已修复 `goroutine/select.go` 中 `break` 只跳出 `select` 导致的死循环和
  `go vet` 不可达代码问题。
- 已修复 `http/main.go` 重复注册 `/` 路由导致服务启动 panic 的问题，改为独立
  `ServeMux`。
- 已抽出 `grpc/client` 客户端求和调用函数，并使用内存 listener 补齐本地 gRPC
  server 调用测试。
- 已完成企业骨架业务可用性：当前用户资料、任务归属隔离、用户与任务分页、稳定业务错误码。
- 已完成企业骨架外部基础设施适配：可选 Redis 用户资料缓存、RabbitMQ 任务事件和实时依赖健康检查；默认关闭，不影响本地运行。
- 已完成企业骨架认证与权限：bcrypt 密码摘要、HS256 JWT、管理员角色控制和内存审计记录。
- 已通过根模块回归：

```powershell
$env:GOMODCACHE="$PWD\.tmp_gomodcache"
$env:GOPATH="$PWD\.tmp_gopath"
go test ./...
```

- 已通过完整检查脚本：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\check.ps1
```

- 已通过包含独立子模块的完整测试脚本：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\test.ps1
```

### 当前注意事项

- 本机 PowerShell 执行策略可能阻止直接运行 `.\scripts\test.ps1`，可使用
  `powershell -ExecutionPolicy Bypass -File .\scripts\test.ps1` 或根模块
  `go test ./...` 兜底。
- `gin.log`、上传文件、trace 文件和本地缓存继续保持忽略，不进入提交。

### 测试覆盖现状

已有测试的目录：

- `algorithm/{binsearch,bubblesort,heapsort,quicksort,selectsort}`
- `array`（当前主要覆盖 `Fibonaci`）
- `designpattern/*`（13 个模式全覆盖）
- `gin`
- `encrypt`、`log`、`runtime`、`task`、`http`
- `goroutine`
- `sync/map`、`sync/mutex`
- `tcp/{client,server}`、`udp/{client,server}`、`grpc/{client,server}`
- `rpc/{arith,client,server}`
- `study`
- 独立子模块 `gomysql`、`goredis`、`gorm`、`gorabbitmq`、`websocket`
- `internal/{auth,repository,router,service}`

完全没有测试或缺少有效默认测试的目录：

- `gmp`

### 当前文档状态

- 常用学习主题、网络示例和 `cmd/server` 均已补充 README。
- 根 `README.md` 已提供二级目录导航，并标注外部服务示例。
- 企业骨架的配置、接口、错误码和可选基础设施说明见 `docs/enterprise-app.md`。

### 代码质量细节

- `array/main.go` 大量示例代码被注释，仅保留 `Fibonaci` 一个可运行入口，与
  “学习示例”定位略有冲突。
- `cmd/server` 的 `internal/response`、`internal/logger` 没有专门测试。

## 第一阶段：基线修正与脚本稳定（已完成）

目标：先修掉最容易影响学习和验证体验的小问题。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| 统一 gRPC 端口示例 | 将根 README 与 `.env.example` 的 `GRPC_ADDR` 默认值对齐 | 已完成 |
| 处理测试脚本执行策略说明 | 在 README 中说明 PowerShell 执行策略限制与替代命令 | 已完成 |
| `scripts/check.ps1` 增加 `go vet` | 提交前同时跑 `go vet ./...` | 已完成 |
| 为 `gmp/` 补 README | 写明 `trace.out` 的生成与查看方式 | 已完成 |
| 清理 `cmd/server/main.go` 空白改动 | 避免无意义 diff 干扰后续提交 | 已完成 |

## 第二阶段：基础示例测试补齐（已完成）

目标：让 `go test ./...` 对默认示例也更有学习价值。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| `encrypt` 测试 | 覆盖 AES + base64 加解密和错误输入 | 已完成 |
| `log` 测试 | 覆盖日志组件和 Gin 日志中间件 | 已完成 |
| `runtime` 测试 | 覆盖 `GOMAXPROCS`、`GOROOT`、`NumCPU` | 已完成 |
| `task` 测试 | 覆盖秒级 cron 表达式和任务 Run | 已完成 |
| `http` 测试 | 使用 `httptest` 覆盖示例 handler 和 mux | 已完成 |
| `goroutine/*` 测试 | 覆盖 WaitGroup、atomic、工作池、select、sync.Map 示例 | 已完成 |
| `sync/{mutex,map}` 测试 | 覆盖并发 map 类型检查和锁示例边界 | 已完成 |
| `tcp/udp/grpc/rpc` 测试 | 客户端服务端使用本地 listener 或服务方法跑通 | 已完成 |
| `study` 测试 | 梳理基础语法示例中可单测的片段 | 已完成 |
| `grpc/client` 测试 | 抽出客户端调用函数后补本地 gRPC server 测试 | 已完成 |
| 独立子模块测试开关 | 外部服务缺失时默认跳过，手动开启时跑真实依赖 | 已完成 |

## 第三阶段：目录文档补齐（已完成）

目标：让新人能根据 README 找到示例、命令和前置条件。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| 根 `README.md` 升级 | 二级目录也列出来；标注哪些需要外部服务 | 已完成 |
| 缺 README 的目录补 README | 为学习专题和网络目录补充运行、测试及前置条件说明 | 已完成 |
| `cmd/server` 模块说明 | 简要介绍 `cmd` 与 `internal` 各包职责 | 已完成 |
| 增加 `CHANGELOG.md` | 记录每个阶段代表性提交 | 已完成 |

## 第四阶段：企业骨架第一阶段（已完成）

目标：把 [enterprise-roadmap.md](enterprise-roadmap.md) 中的“业务可用性”落地。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| 用户资料查询 | 新增 `GET /api/v1/me` | 已完成 |
| 任务归属隔离 | 当前用户只能查看和修改自己的任务 | 已完成 |
| 列表分页 | 用户列表、任务列表支持 `page`、`page_size` | 已完成 |
| 统一错误码 | 业务错误码与 HTTP 状态码分离 | 已完成 |

已完成 roadmap 的“第二阶段：外部基础设施”，当前进入认证与权限阶段。

## 第五阶段：示例能力拓展（可选）

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| `goroutine` 增加 context 示例 | 演示 `context.WithCancel/Timeout` | 单元测试 |
| `gin` 增加中间件组合示例 | Recovery + CORS + 鉴权链 + 限流/超时 | 路由测试 |
| `rpc` 增加 JSON-RPC 或 gRPC gateway 例子 | 在原 TCP RPC 之外补充 | README 给出 curl 例子 |
| 增加 `kafka` 或 `nsq` 例子 | 与 `gorabbitmq` 并列 | 独立 go.mod + README |

## 推荐推进顺序

1. **下一步**：补充 OpenAPI 文档和可选 MySQL 集成测试。
2. **随后**：提供 Docker Compose，编排 MySQL、Redis、RabbitMQ 本地依赖。
3. **后续**：补充 CI 检查示例与审计记录的查询、持久化能力。

## 约束

- 外部服务相关测试默认跳过，通过 `APP_STORAGE`、`APP_CACHE`、`APP_MQ` 等
  环境变量控制是否启用。
- 文档、注释和测试提示继续使用中文。
- 学习和示例目录不引入重型依赖；企业骨架能力优先放在 `cmd/server` 和
  `internal` 下。
- 不提交构建产物、日志、trace、上传文件和本地缓存。
