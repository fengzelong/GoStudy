# udp

这个目录包含基于 UDP 的一次请求、一次回显示例。

| 目录 | 内容 |
| --- | --- |
| `server` | 接收 UDP 数据报并回显给发送方 |
| `client` | 向服务端发送 `Hello server` 并输出响应 |

默认地址为 `127.0.0.1:30000`，可通过 `UDP_ADDR` 修改。

## 运行方式

先启动服务端，再在另一个终端运行客户端：

```powershell
go run ./udp/server
go run ./udp/client
```

## 测试

```sh
go test ./udp/...
```
