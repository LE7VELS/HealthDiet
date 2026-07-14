# Backend

这是项目的 Go + Gin 后端。

## 运行

```powershell
go run ./cmd/api
```

默认地址：`http://localhost:8080`。

可以通过环境变量修改端口：

```powershell
$env:HTTP_ADDR = ":9000"
go run ./cmd/api
```

浏览器访问根路径会返回：

```json
{"message":"HealthDiet API"}
```

## 架构

```text
Router
  ↓
Middleware
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
MongoDB
```

- Router 注册路由。
- Middleware 处理公共 HTTP 逻辑。
- Handler 解析请求和返回响应。
- Service 编写业务逻辑。
- Repository 定义数据访问。
- `repository/mongo` 实现 MongoDB 操作。

依赖只能按上面的方向调用。当前根路径会走完 Router、Middleware、Handler 和 Service；Repository 与 MongoDB 在实现第一个持久化业务时接入。
