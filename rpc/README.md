# rpc

这个目录演示标准库 `net/rpc` 的 HTTP 传输方式。

| 目录 | 内容 |
| --- | --- |
| `arith` | 客户端和服务端共用的乘法、除法服务定义 |
| `server` | 注册 `Arith` 服务并通过 HTTP 提供 RPC |
| `client` | 调用乘法和除法方法 |

默认地址为 `127.0.0.1:8080`，可通过 `RPC_ADDR` 修改。若同时运行 HTTP 或 Gin 示例，请为它们设置不同端口。

## 运行方式

先启动服务端，再在另一个终端运行客户端：

```powershell
go run ./rpc/server
go run ./rpc/client
```

## 测试

```sh
go test ./rpc/...
```
