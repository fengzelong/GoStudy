# gorabbitmq

这个目录演示 RabbitMQ 简单队列模式。公共封装放在当前模块根目录，发送端和接收端拆成两个入口：

- `producer`：发送消息。
- `consumer`：消费消息。

## 前置条件

- 已启动 RabbitMQ 服务。
- 根据本地环境配置 `RABBITMQ_URL` 和 `RABBITMQ_QUEUE`。

示例：

```powershell
$env:RABBITMQ_URL="amqp://guest:guest@127.0.0.1:5672/"
$env:RABBITMQ_QUEUE="GoTestMessage"
go run ./producer
```

另开一个终端启动消费者：

```powershell
$env:RABBITMQ_URL="amqp://guest:guest@127.0.0.1:5672/"
$env:RABBITMQ_QUEUE="GoTestMessage"
go run ./consumer
```

如果没有设置环境变量，程序会使用本地默认示例连接。

