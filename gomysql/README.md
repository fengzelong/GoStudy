# gomysql

这个目录演示使用 `sqlx` 连接 MySQL，并执行新增、查询、更新和删除操作。

## 前置条件

- 已启动 MySQL 服务。
- 已创建示例数据库和 `person` 表。
- 已配置 `MYSQL_DSN`，也可以参考根目录 `.env.example`。

示例：

```powershell
$env:MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/go_test"
go run .
```

如果没有设置 `MYSQL_DSN`，程序会使用本地默认示例连接。

