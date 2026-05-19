# goredis

这个目录演示使用 `redigo` 连接 Redis，并执行设置、获取和批量获取操作。

## 前置条件

- 已启动 Redis 服务。
- 根据本地环境配置 `REDIS_ADDR` 和 `REDIS_PASSWORD`。

示例：

```powershell
$env:REDIS_ADDR="127.0.0.1:6379"
$env:REDIS_PASSWORD=""
go run .
```

如果没有设置 `REDIS_ADDR`，程序会使用 `127.0.0.1:6379`。

