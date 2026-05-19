# websocket

这个目录演示使用 `gorilla/websocket` 和 `gorilla/mux` 实现 WebSocket 服务。

## 运行方式

```powershell
go run .
```

服务默认监听：

```text
127.0.0.1:8080
```

可以通过 `WEBSOCKET_ADDR` 调整监听地址：

```powershell
$env:WEBSOCKET_ADDR="127.0.0.1:8081"
go run .
```

启动后可以打开 `local.html` 或自行连接 `/ws` 路由进行测试。
