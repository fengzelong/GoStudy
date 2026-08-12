# http

这个目录演示标准库 `net/http` 的 `ServeMux`、根路由和用户路径处理。

默认监听地址为 `:8080`，可通过 `HTTP_ADDR` 修改。

## 运行与测试

```sh
go run ./http
go test ./http
```

运行后可访问 `http://127.0.0.1:8080/` 和 `http://127.0.0.1:8080/user/Alice`。
