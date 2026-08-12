# sync

这个目录演示 `sync` 包的常用并发同步工具。

| 目录 | 内容 |
| --- | --- |
| `map` | 在 `sync.Map` 外增加键和值类型校验 |
| `mutex` | 通过读写黑板示例演示 `sync.RWMutex` |

## 运行与测试

```sh
go test ./sync/...
go run ./sync/map
go run ./sync/mutex
```

`mutex` 示例包含多次 `time.Sleep`，运行时会持续输出读写过程。
