# Backend

这是项目的 Go + Gin 后端。

## 运行

```powershell
go run .
```

默认地址：`http://localhost:8080`。

可以通过环境变量修改端口：

```powershell
$env:HTTP_ADDR = ":9000"
go run .
```

浏览器访问根路径会返回：

```json
{"message":"HealthDiet API"}
```

## 架构

```text
Router + Middleware
          ↓
       Handler
       ↙     ↘
 简单 CRUD   复杂业务
     ↓          ↓
   Store     Service
     ↓          ↓
     └──── Store
             ↓
          MongoDB
```

- Router 注册路由。
- Middleware 处理公共 HTTP 逻辑。
- Handler 解析请求和返回响应。
- 简单 CRUD 可以由 Handler 直接调用 Store。
- Service 只编写包含规则、计算或多个步骤的复杂业务。
- Model 保存 Handler、Service 和 Store 共享的业务结构，不是必须经过的一层。
- Store 集中实现 MongoDB 操作。

当前根路径只经过 Router、Middleware 和 Handler。Store 在实现第一个持久化业务时接入，Service 在出现真实复杂业务时按需加入。
