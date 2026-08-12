# tcp

这个目录包含基于 TCP 的回显服务端和交互式客户端。

| 目录 | 内容 |
| --- | --- |
| `server` | 接收客户端数据并原样返回 |
| `client` | 从标准输入读取一行数据并发送；输入 `Q` 退出 |

默认监听和连接地址为 `127.0.0.1:20000`，可通过 `TCP_ADDR` 修改。

## 运行方式

先启动服务端：

```powershell
go run ./tcp/server
```

在另一个终端启动客户端：

```powershell
go run ./tcp/client
```

## 测试

```sh
go test ./tcp/...
```
