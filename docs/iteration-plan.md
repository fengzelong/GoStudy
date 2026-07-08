# 项目迭代计划

这份规划覆盖 GoStudy 仓库里尚未整理或尚未补齐的部分，包含 `gin/` 模块当前
未提交的改动、基础示例的测试与文档缺口，以及企业骨架的下一阶段任务。它与
[enterprise-roadmap.md](enterprise-roadmap.md) 互补：本文侧重全局清理与基础
示例，roadmap 文档继续负责企业骨架的纵向推进。

## 现状盘点

### 工作区未提交改动

`gin/` 目录有较大改动，尚未提交：

| 文件 | 状态 | 摘要 |
| --- | --- | --- |
| `gin/main.go` | 已修改 | 拆出 `RegisterBasicRoutes` 等多个注册函数；新增 `/ping`、`/param`、`/query`、`/header`、`/cookie/*`、`/bindQuery`、`/form`、`/bindHeader`、`/bindXML`、`/download`、`/stream`、`/basic/secret`；增加 `Recovery` 优雅关闭、日志文件路径、NoRoute/NoMethod |
| `gin/README.md` | 新增 | 列出环境变量、路由清单、curl 示例、测试命令 |
| `gin/main_test.go` | 新增 | 4 个测试函数覆盖 ping、绑定、表单/Header/Cookie/Auth、上传 |
| `.env.example` | 已修改 | 新增 `GIN_LOG_FILE`、`GIN_UPLOAD_DIR` |

这一组改动需要在「第一阶段」开始前先收尾，确保默认测试可运行。

### 测试覆盖现状

已有测试的目录：

- `algorithm/{binsearch,bubblesort,heapsort,quicksort,selectsort}`
- `array`（仅 `Fibonaci`）
- `designpattern/*`（13 个模式全覆盖）
- `internal/{auth,repository,router,service}`
- `rpc/arith`
- `gin`（新增，待提交）

完全没有测试的目录：

- `encrypt`、`log`、`runtime`、`study`、`task`、`http`
- `goroutine/*`（atomic、channel、map、pool、select、waitgroup 共 6 个文件）
- `sync/mutex`、`sync/map`
- `gmp`（仅 `trace.go`）
- `tcp/{client,server}`、`udp/{client,server}`、`grpc/{client,server}`、`rpc/{client,server}`
- 独立子模块 `gomysql`、`goredis`、`gorm`、`gorabbitmq`、`websocket`

### 文档缺口

- 一级目录以下没有 `README`：`algorithm`、`array`、`designpattern`、`encrypt`、
  `log`、`runtime`、`study`、`task`、`http`、`goroutine`、`sync`、`gmp`、
  `tcp`、`udp`、`grpc`、`proto`、`rpc`。
- `cmd/server` 与 `internal/*` 没有模块级 README，仅靠根 `docs/enterprise-app.md`
  说明。
- 根 `README.md` 的目录说明只到一级目录，没有体现子目录结构。

### 代码质量细节

- `array/main.go` 大量示例代码被注释，仅保留 `Fibonaci` 一个可运行入口，与
  「学习示例」定位略有冲突。
- `gmp/trace.go` 写死 `trace.out`，没有任何说明如何用 `go tool trace` 观察。
- `cmd/server` 的 `internal/response`、`internal/logger` 没有专门测试。
- `scripts/check.ps1` 只检查 git 中变动的 `.go` 文件，缺少 `go vet` 步骤。
- `docs/enterprise-roadmap.md` 中 GRPC 端口示例写的是 `:50053`，与代码及
  `.env.example` 中的 `:50052` 不一致。

## 第一阶段：工作树收尾与基线

目标：把现有未提交改动验证落地，修复最容易观察到的工程小坑。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| 收尾 `gin/` 改动 | 在本机运行 `go run ./gin` 冒烟、`go test ./gin` 通过 | `scripts/test.ps1` 全部通过；提交记录说明新增能力 |
| 修正端口示例 | `docs/enterprise-roadmap.md` 中 GRPC 端口统一为 `:50052` | grep 确认无残留 `:50053` |
| `scripts/check.ps1` 增加 `go vet` | 提交前同时跑 `go vet ./...` | 故意引入可疑代码时脚本应失败 |
| 整理 `array/` | 把 `main.go` 中被注释的示例拆成单独子目录或保持单文件但补注释 | 每个示例都有 `// 示例名：xxx` 注释 |
| 为 `gmp/` 补 README | 写明 `trace.out` 的查看方式 | README 中给出 `go tool trace` 命令 |

## 第二阶段：基础示例测试补齐

目标：让 `go test ./...` 对默认示例也有意义。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| `encrypt` 测试 | 覆盖常用哈希、对称加密 | `go test ./encrypt` |
| `log` 测试 | zap 与 lumberjack 切割 | `go test ./log` |
| `runtime` 测试 | 覆盖 `runtime.GOROOT`、`NumCPU` 等基础用法 | `go test ./runtime` |
| `task` 测试 | 简单 task 调度 | `go test ./task` |
| `http` 测试 | httptest 覆盖示例 server | `go test ./http` |
| `goroutine/*` 测试 | channel、select、WaitGroup、atomic、map、pool 各自可单测部分 | `go test ./goroutine` |
| `sync/{mutex,map}` 测试 | 锁与并发 map 行为 | `go test ./sync/...` |
| `tcp/udp/grpc/rpc` 测试 | 客户端服务端起本地 listener 跑通 | `go test ./tcp/... ./udp/... ./grpc/... ./rpc/...` |
| 独立子模块测试开关 | `scripts/test.ps1` 中各子模块 `go test ./...` 的预期结果（外部服务缺失时允许跳过） | README 给出运行前置 |

## 第三阶段：文档与目录说明

目标：让新人能根据 README 找到示例和命令。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| 根 `README.md` 升级 | 二级目录也列出来；标注哪些需要外部服务 | 表格里能看到 `algorithm/binsearch` 等 |
| 缺 README 的目录补 README | 至少为 `algorithm`、`array`、`designpattern`、`goroutine`、`sync`、`gmp`、`tcp`、`udp`、`grpc`、`proto`、`rpc` 各写一份 | `find . -maxdepth 3 -name README.md` 数量显著增加 |
| `cmd/server` 模块说明 | 简要介绍 `cmd` 与 `internal` 各包职责 | 链接到 `docs/enterprise-app.md` |
| 增加 `CHANGELOG.md` | 记录每个阶段的代表性提交 | 至少包含本次 `gin/` 增强 |

## 第四阶段：企业骨架第一阶段（沿用 roadmap）

目标：把 [enterprise-roadmap.md](enterprise-roadmap.md) 中的「业务可用性」
落地。

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| 用户资料查询 | 新增 `GET /api/v1/me` | 路由测试覆盖未登录、已登录 |
| 任务归属隔离 | 跨用户访问被拒绝 | 路由测试断言 403/404 |
| 列表分页 | 用户、任务支持 `page`、`page_size` | 服务层 + 路由层单测 |
| 统一错误码 | 业务错误码与 HTTP 状态码分离 | 文档列出错误码表，测试断言响应结构 |

完成后即可进入 roadmap 的「第二阶段：外部基础设施」。

## 第五阶段：示例能力拓展（可选）

| 任务 | 说明 | 验收方式 |
| --- | --- | --- |
| `goroutine` 增加 context 示例 | `context.WithCancel/Timeout` 演示 | 单元测试 |
| `gin` 增加 middleware 组合示例 | Recovery + CORS + 鉴权链 | 路由测试 |
| `rpc` 增加 JSON-RPC 或 gRPC gateway 例子 | 在原 TCP RPC 之外补充 | README 给出 curl 例子 |
| 增加 `kafka` 或 `nsq` 例子 | 与 `gorabbitmq` 并列 | 独立 go.mod + README |

## 推荐推进顺序

1. **立刻**：跑 `scripts/test.ps1` 验证 `gin/` 改动，修复 roadmap 端口笔误，
   补 `gmp` README。
2. **1 周内**：把 `encrypt`、`log`、`runtime`、`task`、`http` 这些小模块的
   测试补齐。
3. **1~2 周**：补齐 `goroutine/*`、`sync/*` 测试和 `tcp/udp/grpc/rpc` 测试。
4. **2~3 周**：写各目录 README，升级根 `README.md`。
5. **3 周后**：进入企业骨架第一阶段（`/me`、任务归属、分页、错误码）。

## 约束

- 现有未提交改动属于工作区内容，不在本次计划里「回滚或重写」。
- 外部服务相关测试默认跳过，通过 `APP_STORAGE`、`APP_CACHE`、`APP_MQ` 等
  环境变量控制是否启用。
- 文档、注释、测试提示继续使用中文。
- 学习和示例目录不引入重型依赖；新能力优先放在 `cmd/server` 和 `internal`
  下。
