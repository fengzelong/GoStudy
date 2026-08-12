# goroutine

这个目录演示 Go 并发基础：goroutine、channel、`select`、`WaitGroup`、读写锁、工作池、`sync.Map` 和原子操作。

主要示例包括 channel、工作池、多路 channel 选择、并发安全映射和原子加法。

## 运行与测试

```sh
go run ./goroutine
go test ./goroutine
```

`main.go` 默认运行原子加法示例；其余示例可按需取消注释。部分示例会创建大量 goroutine 或带有等待时间。
