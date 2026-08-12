# task

这个目录演示 `robfig/cron` 定时任务。示例使用支持秒字段的 cron 解析器，并每五秒执行一次任务。

## 运行与测试

```sh
go run ./task
go test ./task
```

示例会持续运行，使用 `Ctrl+C` 停止。
