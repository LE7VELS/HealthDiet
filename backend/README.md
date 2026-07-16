# Backend

这是项目的 Go + Gin 后端。

## 运行

```powershell
go run .
```

启动前需要本机 MongoDB 可用。默认配置：

```text
HTTP_ADDR=:8080
MONGODB_URI=mongodb://127.0.0.1:27017
MONGODB_DATABASE=healthdiet
JWT_SECRET=至少32个字符的随机密钥
CORS_ALLOWED_ORIGIN=http://localhost:5173
```

本地运行会自动读取后端目录中的 `.env`，系统环境变量优先于文件配置。`JWT_SECRET` 没有默认值，未设置或少于 32 个字符时服务会拒绝启动，避免使用写在代码里的固定密钥；仓库只提交 `.env.example`，本地 `.env` 已被 Git 忽略。

认证接口位于 `/api/v1/auth`：注册和登录会返回一小时有效的 Bearer JWT，`GET /api/v1/auth/me` 用于验证登录状态。前端本地开发默认通过 Vite 把 `/api` 请求代理到 `http://127.0.0.1:8080`。

默认 API 地址为 `http://localhost:8080`。可以通过环境变量覆盖上述配置，例如：

```powershell
$env:JWT_SECRET = "请替换为至少32个字符的本地随机密钥"
$env:HTTP_ADDR = ":9000"
go run .
```

启动时后端会连接并 Ping MongoDB，然后幂等创建 [`../DATA_MODEL.md`](../DATA_MODEL.md) 定义的当前集合和索引。连接或初始化失败时，Gin 不会启动。

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

- Router 的 `router.go` 创建 Engine 和全局 Middleware，`root_routes.go`、`auth_routes.go` 等文件按已实现业务模块注册路由。
- Middleware 处理公共 HTTP 逻辑。
- Handler 解析请求和返回响应。
- 简单 CRUD 可以由 Handler 直接调用 Store。
- Service 只编写包含规则、计算或多个步骤的复杂业务。
- Model 保存 Handler、Service 和 Store 共享的业务结构，不是必须经过的一层。
- Store 集中实现 MongoDB 连接、集合与索引初始化以及后续数据访问。

当前 `main.go` 创建并管理共享 Store 的生命周期，再依次装配认证 Service、Handler 和 `router.Dependencies`。Router 的 `New` 只接收这个带命名字段的依赖结构体，后续增加用户资料、Food 等 Handler 时不会形成很长的位置参数列表；Store 仍只注入 Handler 或 Service，不直接交给 Router。

当前路由目录：

```text
internal/router/
├─ router.go          # Engine、全局 Middleware 和 Dependencies
├─ root_routes.go     # 根路径和统一 404
└─ auth_routes.go     # 注册、登录、当前用户
```

用户资料或 Food 接口真正实现后，再按相同方式增加 `user_routes.go`、`food_routes.go`，不提前保留空文件。
