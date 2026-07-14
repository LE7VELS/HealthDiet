# 后端开发约束

> 功能见 [`BACKEND_REQUIREMENTS.md`](./BACKEND_REQUIREMENTS.md)，接口见 [`../API_CONTRACT.md`](../API_CONTRACT.md)。

## 1. 技术栈

- Go + Gin。
- MongoDB 官方 Go Driver v2。
- JWT。
- Docker Compose。

不使用 GORM，不增加第二个后端服务、第二种数据库或复杂基础设施。

## 2. 简单目录结构

```text
backend/
├─ cmd/api/                 # 程序入口
├─ internal/
│  ├─ config/               # 环境变量
│  ├─ router/               # Router
│  ├─ middleware/           # Middleware
│  ├─ handler/              # Handler 和 HTTP DTO
│  ├─ service/              # Service
│  └─ repository/
│     └─ mongo/             # Repository 的 MongoDB 实现
├─ go.mod
└─ go.sum
```

业务增多后可以在 Handler、Service 和 Repository 内按 `auth`、`food`、`meal` 等模块拆文件，但不额外增加一套分层方式。

## 3. 职责

- `main.go`：读取配置，从 Repository 开始向上组装依赖并启动 Gin。
- Router：注册路径和 Middleware，只调用 Handler。
- Middleware：处理所有路由共享的 HTTP 逻辑，不承载业务规则。
- Handler：解析请求、调用 Service、返回响应。
- Service：业务规则、所有权校验和营养计算流程，只依赖 Repository 接口。
- Repository：定义数据访问接口；`repository/mongo` 实现 MongoDB 查询和保存。

依赖方向固定为：

```text
Router → Middleware → Handler → Service → Repository → MongoDB
```

- Router 不直接调用 Service 或 Repository。
- Middleware 不调用 Repository。
- Handler 不直接操作 MongoDB。
- Service 不依赖 Gin Context，也不导入 MongoDB Driver。
- MongoDB BSON、Collection 和 Driver 类型只出现在 `repository/mongo` 与必要的数据转换边界。

## 4. Go 规则

- 使用 `gofmt`。
- 显式处理错误，避免无意义的 `panic`。
- 保持函数和包职责清楚，不为了“分层”制造大量只有一行的文件。
- 配置集中在 `internal/config`，密钥和连接地址来自环境变量。
- 新依赖必须有明确用途。

## 5. MongoDB

- 只使用官方 Go Driver v2。
- 复用一个 MongoDB Client，并给数据库操作传递 `context.Context`。
- 私有数据查询必须包含 `user_id`。
- 列表查询限制数量并分页。
- 用户名、邮箱和必要一对一关系使用唯一索引。
- MongoDB 错误在 Repository 或 Service 转换，不直接返回客户端。

## 6. 认证和上传

- 密码使用安全单向哈希，不记录密码或完整 JWT。
- JWT 必须验证签名和有效期。
- 所有者来自认证上下文，不信任请求中的 `userId`。
- 图片检查实际类型、大小和所有权。
- 文件名使用服务端生成的存储标识，防止路径穿越和重名覆盖。

## 7. 营养计算

- 只使用 Food、用户确认的重量或份量和确定性公式。
- 前端提交的合计不可信。
- 记录保存时生成营养快照。
- 缺失值和真实零值必须区分。
- 每日汇总和 7 天趋势复用同一计算规则。

## 8. 检查方式

当前以手工功能测试为主，不要求先建设完整自动化测试体系。

每次修改后至少执行：

```powershell
gofmt -w cmd internal
go vet ./...
go build ./cmd/api
```

## 9. Git 和安全

- 不提交真实 `.env`、密钥、数据库数据、上传文件、日志和构建产物。
- 不覆盖无关改动。
- 未经授权不执行 commit、push、rebase、reset 或强制操作。
