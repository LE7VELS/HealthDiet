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
├─ main.go                  # 单个 API 程序入口
├─ internal/
│  ├─ config/               # 环境变量
│  ├─ router/               # Router
│  ├─ middleware/           # Middleware
│  ├─ handler/              # Handler 和 HTTP DTO
│  ├─ service/              # 复杂业务 Service，按需使用
│  ├─ model/                # 共享业务 Model
│  └─ store/                # MongoDB 数据访问
├─ go.mod
└─ go.sum
```

当前只有一个 API 可执行程序，因此入口直接放在后端根目录。以后出现第二个可执行程序或后台任务时，再按需要恢复 `cmd/<name>`。

业务增多后可以在 Handler、Service、Model 和 Store 内按 `auth`、`food`、`meal` 等模块拆文件，但不额外增加一套分层方式。

## 3. 职责

- `main.go`：读取配置、创建共享依赖并启动 Gin。
- Router：注册路径和 Middleware，只调用 Handler。
- Middleware：处理所有路由共享的 HTTP 逻辑，不承载业务规则。
- Handler：解析请求和返回响应；简单的单集合 CRUD 可以直接调用 Store。
- Service：只负责包含业务规则、多个步骤、跨集合操作、状态变化或营养计算的流程。
- Model：定义 User、Profile、Food、Meal、Nutrients 等共享业务结构，不包含 HTTP 或 MongoDB 实现细节。
- Store：集中实现 MongoDB 查询和保存，并要求私有查询显式接收当前用户 ID。

允许的依赖方向为：

```text
Router + Middleware → Handler → Store → MongoDB
Router + Middleware → Handler → Service → Store → MongoDB
```

- Router 不直接调用 Service 或 Store。
- Middleware 不调用 Store。
- Handler 可以调用 Store，但不能导入 MongoDB Driver、BSON 或 Collection。
- Service 不依赖 Gin Context，也不导入 MongoDB Driver。
- Handler 的请求和响应 DTO 放在对应 Handler 文件中；Service 输入放在对应 Service 文件中，不集中创建庞大的 DTO 包。
- Model 不依赖 Gin，不包含 JSON 请求校验、HTTP 状态码、MongoDB `ObjectID` 或 BSON 标签。
- MongoDB BSON、Collection 和 Driver 类型只出现在 Store 与必要的数据转换边界。
- Store 内部可以定义未导出的 MongoDB Document，并负责在 Document 与 Model 之间转换。
- 不要求每个接口都创建 Service；只做参数转发的 Service 应省略。
- Store 先使用具体实现；只有出现测试替身、第二种实现或明确解耦需求时，才在使用方定义小接口。

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
- MongoDB 错误在 Store 或 Service 转换，不直接返回客户端。

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
go fmt ./...
go vet ./...
go build .
```

## 9. Git 和安全

- 不提交真实 `.env`、密钥、数据库数据、上传文件、日志和构建产物。
- 不覆盖无关改动。
- 未经授权不执行 commit、push、rebase、reset 或强制操作。

## 10. 后续 AI 接入约束

以下规则只约束未来扩展，当前不要求创建 AI 相关代码或目录：

- AI 不直接连接包含用户私有数据的核心业务库，用户上下文由 Go API 按认证用户范围提供。
- 可以为隔离后的公共食品知识库或只读视图创建独立的最小权限只读账号；该账号不能访问核心业务集合，也不能拥有写权限。
- Agent 工具 Handler 仍调用现有 Service，不复制营养计算、所有权校验或写入规则。
- 可以提供聚合上下文工具减少 Agent 调用次数，但不能接受由模型提交的任意 MongoDB 查询。
- 所有 AI 写操作都由 Go Service 校验，并在需要时要求用户确认。
