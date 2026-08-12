# grpc

这个目录演示 gRPC 服务端与客户端。协议定义和生成代码位于 [`proto/grpc`](../proto/grpc)。

服务提供 `SayHello` 和 `SumFunc` 两个方法，默认服务地址为 `127.0.0.1:50052`，可通过 `GRPC_ADDR` 修改。运行时 trace 页面默认监听 `:50051`，可通过 `GRPC_TRACE_ADDR` 修改。

## 运行方式

先启动服务端：

```powershell
go run ./grpc/server
```

再在另一个终端运行客户端：

```powershell
go run ./grpc/client
```

## 测试

```sh
go test ./grpc/...
```

客户端测试会使用内存 listener，不需要手工启动服务端。
