# log

这个目录演示 zap 结构化日志、lumberjack 日志切割，以及 Gin 日志中间件。

运行后会启动一个监听 `:8080` 的 Gin 服务，并在当前目录写入 `zap.log`。

## 运行与测试

```sh
go run ./log
go test ./log
```

`zap.log` 是运行时产物，不需要提交。
