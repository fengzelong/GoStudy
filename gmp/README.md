# gmp

这个目录演示 Go 运行时 trace 的最小用法：程序启动后创建 `trace.out`，通过
`runtime/trace` 记录运行时事件。

## 运行方式

```powershell
go run ./gmp
```

运行后会在当前目录生成 `trace.out`。查看 trace：

```powershell
go tool trace trace.out
```

`trace.out` 是运行时产物，已经由 Git 忽略，不需要提交。
