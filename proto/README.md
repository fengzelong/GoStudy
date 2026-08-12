# proto

这个目录保存协议定义和生成代码。目前包含 `grpc/grpc.proto`，它定义了示例 gRPC 服务的消息和接口；`grpc.pb.go` 是对应的 Go 生成代码，由 `grpc` 客户端和服务端共用。

## 使用方式

```sh
go test ./grpc/...
```

修改 `.proto` 后，需要使用与项目依赖兼容的 `protoc` 和 Go 插件重新生成 `grpc.pb.go`，并将协议与生成代码一并提交。
