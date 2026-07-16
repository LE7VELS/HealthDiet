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
│  ├─ router/               # Router 装配和按业务模块拆分的路由注册
│  │  ├─ router.go          # Engine、全局 Middleware、Dependencies
│  │  ├─ root_routes.go     # 根路径和统一 404
│  │  └─ auth_routes.go     # 当前已实现的认证路由
│  ├─ middleware/           # Middleware
│  ├─ handler/              # Handler 和 HTTP DTO
│  ├─ service/              # 复杂业务 Service，按需使用
│  ├─ model/                # 共享业务 Model
│  └─ store/                # MongoDB 数据访问
├─ go.mod
└─ go.sum
```

当前只有一个 API 可执行程序，因此入口直接放在后端根目录。以后出现第二个可执行程序或后台任务时，再按需要恢复 `cmd/<name>`。

Router、Handler、Service、Model 和 Store 可以按已经实现的 `auth`、`food`、`meal` 等业务模块拆文件，但不额外增加一套分层方式。路由文件使用 `<module>_routes.go` 命名；只有模块出现真实接口时才创建对应文件，不提前增加空的 `user_routes.go`、`food_routes.go` 等占位文件。

## 3. 职责

- `main.go`：读取配置、创建共享依赖并启动 Gin。
- Router：`router.go` 创建 Engine、注册全局 Middleware 并依次调用模块路由注册函数；`*_routes.go` 只注册本模块路径、路由级 Middleware 和 Handler。
- Middleware：处理所有路由共享的 HTTP 逻辑，不承载业务规则。
- Handler：解析请求和返回响应；简单的单集合 CRUD 可以直接调用 Store。
- Service：只负责包含业务规则、多个步骤、跨集合操作、状态变化或营养计算的流程。
- Model：定义 User、Profile、Food、Meal、Nutrients 等共享业务结构，不包含 HTTP 或 MongoDB 实现细节。
- 食品和菜谱共用 Food Model、Handler 和 Store，通过 `kind` 和可选菜谱内容区分；当前不拆分独立 Recipe 分层。
- Store：集中实现 MongoDB 查询和保存，并要求私有查询显式接收当前用户 ID。

`main.go` 创建一个应用共享的 Store，并负责在进程退出时关闭 MongoDB Client；同时按依赖方向创建 Service 和 Handler。Router 的 `New` 使用一个带命名字段的 `Dependencies` 结构体接收已构造的 Handler、Middleware 所需组件和路由配置，避免模块增加后出现很长的位置参数列表。该结构体不能成为收纳任意全局对象的“万能容器”，只添加路由注册实际使用的依赖，也不能把 Store 直接传给 Router。

允许的依赖方向为：

```text
Router + Middleware → Handler → Store → MongoDB
Router + Middleware → Handler → Service → Store → MongoDB
```

- Router 不直接调用 Service 或 Store。
- Router 的模块注册函数保持未导出，只由 `router.go` 统一编排；新增模块时在 `New` 中显式注册，避免依赖文件初始化顺序或隐藏的全局状态。
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
- `fmt.Errorf("操作说明: %w", err)` 用于在 Store 和 Service 中补充错误上下文；错误到达 Handler 等 HTTP 边界时必须有日志出口，不能只返回通用响应后丢弃详细错误。
- 预期业务失败记录稳定错误码、人类可读安全消息、路径和必要请求元数据即可；未知内部错误记录完整包装错误链。禁止记录明文密码、完整 JWT、连接串和敏感请求体，也不通过日志区分“账号不存在”和“密码错误”。
- 保持函数和包职责清楚，不为了“分层”制造大量只有一行的文件。
- 后端代码需要包含充分且准确的中文注释，不能只在文件顶部或入口函数放一条笼统说明。
- 每个包要说明职责与依赖边界；导出的类型和函数必须说明用途。未导出的 HTTP DTO、MongoDB Document、错误类型和安全辅助函数，只要字段含义、转换方向或使用限制不直观，也必须说明。
- Handler 注释要覆盖请求 DTO、响应 DTO、统一错误映射和不可信输入边界；Service 注释要覆盖业务步骤、规范化规则、并发冲突兜底和敏感数据边界；Store 注释要覆盖 BSON 与 Model 转换、查询范围、唯一索引及 Driver 错误转换；Middleware 注释要覆盖上下文写入、拒绝条件和跨域安全策略。
- 复杂流程应在关键分支前说明“为什么”，例如 JWT 验签顺序、密码哈希边界、并发注册依赖唯一索引、未知数据库错误不得返回客户端。简单赋值和一眼可懂的控制流不要求逐行注释，避免无信息量注释。
- 完成后端修改前必须逐个检查本次涉及的 Go 文件，确认包职责、主要类型、导出函数、关键私有辅助函数和安全边界已有注释；代码变化时同步维护，不保留“后续实现”等与现状不符的旧注释。
- 配置集中在 `internal/config`，密钥和连接地址来自环境变量。
- 新依赖必须有明确用途。

## 5. MongoDB

- 只使用官方 Go Driver v2。
- 复用一个 MongoDB Client，并给数据库操作传递 `context.Context`。
- MongoDB 连接地址和数据库名分别来自 `MONGODB_URI` 和 `MONGODB_DATABASE`，本地默认值为 `mongodb://127.0.0.1:27017` 和 `healthdiet`。
- 应用启动时 Ping MongoDB，并幂等创建 `DATA_MODEL.md` 约定的集合和索引；初始化失败时不启动 Gin。
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
- 菜谱创建和修改包含所有权校验、多个 Food 查询、快照与营养计算，必须由 Service 编排，Handler 不得直接拼装计算结果。
- 菜谱原料只允许引用当前用户可见且 `kind=food` 的 Food；暂不支持嵌套菜谱。
- 菜谱总营养是各原料营养之和；每份营养按份数计算；每 100g 营养优先按成品重量计算，否则按原料总重量估算并明确记录依据。
- 制作步骤不参与当前公式。不得通过关键词或自然语言猜测烹饪损耗、吸油量或成品重量。
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
