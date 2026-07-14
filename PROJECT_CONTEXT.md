# 项目背景与开发上下文

> 本文档说明项目当前阶段、范围和技术边界。开始任务前先阅读本文档。

## 1. 项目背景

这是一个面向个人用户的饮食记录与营养分析系统，帮助用户记录饮食、查看每日营养摄入和最近趋势。

本项目只提供饮食参考，不提供医疗诊断或治疗建议。

## 2. 产品主链路

```text
用户注册登录
  → 建立个人档案和营养目标
  → 查询或创建食物
  → 添加饮食记录
  → 保存到 MongoDB
  → 查看每日营养汇总
  → 查看最近 7 天趋势
```

## 3. 当前开发阶段

当前同时开发 React 前端和 Go 后端，目标是完成一个本地可真实运行的小型系统。

- 使用真实 Go API，不以 Mock 演示作为完成标准。
- 数据保存到 MongoDB。
- 当前以手工功能测试为主。
- 开发时保持结构清晰，但不提前建设复杂基础设施。

## 4. 当前前端范围

- 登录和注册。
- Dashboard。
- 饮食记录。
- 食物搜索、详情和自定义食物。
- 每日营养汇总和最近 7 天趋势。
- 营养目标、饮食偏好和个人档案。
- 图片附件选择和上传。
- 响应式电脑端和手机端。

## 5. 当前后端范围

- Go + Gin API。
- JWT 注册登录和认证。
- 用户档案、营养目标和饮食偏好。
- Food 查询和自定义 Food。
- 饮食记录 CRUD。
- 图片上传和元数据。
- 每日营养汇总和最近 7 天趋势。
- MongoDB 持久化。
- Docker Compose 本地环境。

## 6. 当前食物数据方案

Food 数据只来自：

1. 项目 Seed 数据。
2. JSON 或 CSV 开发导入。
3. 用户创建的自定义 Food。

找不到食物时，用户可以创建自定义 Food，不能因为数据不全阻塞主流程。

## 7. 当前明确不做

- AI Agent、LLM、Python Agent、Prompt、Embedding 和向量数据库。
- 自动食谱生成和完整个性化推荐。
- 食品网站爬虫和公共食品 API。
- 图片识别食物和自动估算重量。
- 自然语言录入。
- 原生移动 App。
- 医疗诊断和治疗建议。
- 30 天、自定义区间等长期分析。

## 8. 技术栈

### 前端

- React + TypeScript + Vite。
- Astryx Design System + StyleX。
- React Router。
- TanStack Query。
- React Hook Form + Zod。
- ECharts。

### 后端

- Go + Gin。
- MongoDB 官方 Go Driver v2。
- JWT。
- Docker Compose。

MongoDB 不使用 GORM。图片由 Go 服务处理，MongoDB 只保存图片元数据和访问标识。

## 9. 当前整体架构

```text
React Web
   │ HTTP / JWT / 图片
   ▼
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

- React 只调用 Go API。
- Router 只注册路径和中间件。
- Middleware 处理日志、异常恢复和后续 JWT 认证。
- Handler 只处理 HTTP 输入输出。
- Service 负责业务规则和营养计算流程。
- Repository 负责数据访问，MongoDB 细节只存在于 Repository 实现中。
- 依赖只能按上面的方向向下调用，不能跨层或反向依赖。
- 营养值来自 Food、用户确认的重量或份量和确定性公式。

## 10. 食物数据边界

- 当前不追求完整互联网食品库。
- 不猜测缺失营养值，缺失和真实零值必须区分。
- Food 修改不能悄悄改变历史饮食记录，记录保存时保留营养快照。
- 爬虫和公共食品 API 留到后续阶段。

## 11. 后续阶段

1. 扩展 Food 数据来源。
2. 增加更长时间的历史分析和更多规则提示。
3. 数据与规则稳定后再评估 AI Agent。

## 12. AI 协作规则

未来 AI 只能读取 Go API 的结构化数据、调用受控工具和解释后端结果。AI 不得直接访问 MongoDB、计算营养值、生成虚假 Food 或绕过 Go API 写数据。

开发协作时：

- 不扩大当前任务范围。
- 修改前检查 Git 状态，保留用户已有改动。
- 不提交密钥、真实 `.env`、依赖、构建产物和本地上传文件。
- 未经授权不执行 commit、push、rebase、reset 或强制操作。

## 13. 文档职责

- `PROJECT_CONTEXT.md`：项目阶段和整体边界。
- `frontend/FRONTEND_REQUIREMENTS.md`：前端页面和交互。
- `frontend/FRONTEND_CONSTRAINTS.md`：前端实现约束。
- `backend/BACKEND_REQUIREMENTS.md`：后端功能范围。
- `backend/BACKEND_CONSTRAINTS.md`：后端实现约束。
- `API_CONTRACT.md`：前后端接口合同。
- `DATA_MODEL.md`：MongoDB 数据模型。

文档冲突时先指出，不自行猜测。
