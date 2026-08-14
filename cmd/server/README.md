# 企业应用服务入口

`cmd/server` 是企业应用骨架的可执行入口。它读取环境变量、初始化应用和日志，并在收到终止信号后优雅关闭 HTTP 服务。

业务实现放在 [`internal`](../../internal) 中，按应用装配、路由、服务、仓储、领域对象、鉴权和中间件分层，避免把业务代码堆放在入口函数。

## 运行方式

默认使用内存仓储，缓存与消息队列关闭，不依赖外部服务：

```sh
go run ./cmd/server
```

启动后可访问：

```sh
curl http://127.0.0.1:8080/health
```

配置项、接口清单、MySQL 切换方法和分层说明见 [企业应用骨架说明](../../docs/enterprise-app.md)；后续任务见 [企业应用迭代规划](../../docs/enterprise-roadmap.md)。

当前服务已支持当前用户资料、分页列表、任务归属隔离、稳定业务错误码，以及可选的用户资料缓存和任务事件发布。可通过 `/health` 查看存储、缓存和消息队列依赖状态。

接口契约见 [OpenAPI 文档](../../docs/openapi.yaml)。如需使用 MySQL、Redis、RabbitMQ，可在仓库根目录执行 `docker compose up -d` 启动本地依赖；应用本身仍在本机以 `go run ./cmd/server` 运行。

## 测试

入口本身没有单独测试，相关路由、服务、仓储和鉴权测试可通过以下命令运行：

```sh
go test ./internal/...
```
