# 前端开发约束

## 1. 当前范围

- 当前只实现前端，不引入 AI，不编写真实后端逻辑。
- 产品为响应式 Web，同时支持电脑和手机。
- 优先完成：登录注册、个人档案、食品库、饮食记录、营养报告和推荐页面。

## 2. 技术栈

- React + TypeScript + Vite。
- UI 使用 Astryx Design System，自定义样式优先使用其 Token 和 StyleX。
- 路由使用 React Router，服务端数据状态使用 TanStack Query。
- 表单使用 React Hook Form + Zod，图表使用 ECharts。
- 不引入 Redux，除非后续出现明确的复杂全局状态。

## 3. 开发规则

- 使用 TypeScript strict 模式，避免 `any`。
- 业务按 `features` 组织；Astryx 组件统一经 `components/ui` 封装后使用。
- API 请求统一放在 `lib/api`，地址来自 `VITE_API_BASE_URL`，禁止在页面内硬编码。
- 后端完成前使用独立 Mock 数据，页面不得与 Mock 实现强耦合。
- 上传组件必须支持：电脑选择/拖拽文件，手机拍照/相册选择，并提供预览、类型和大小校验。
- 图片上传必须通过 `lib/api` 的统一接口；未来接入真实服务时，由 Go API 负责本地文件系统或 OSS 等对象存储的持久化，前端不得硬编码 OSS 密钥和真实存储地址。
- 页面必须同时在常见手机宽度和桌面宽度下可用，不允许只做桌面端后简单缩放。

## 4. Git 约束

- 项目根目录必须保留 `.gitignore`。
- 禁止提交 `node_modules`、`dist`、`.env`、日志、缓存、IDE 配置和本地上传文件。OSS 接入代码与 `.env.example` 等配置示例可以提交，但不得包含真实密钥、Token 或服务地址。
- 密钥、Token 和真实服务地址不得入库；只提交 `.env.example`。
- 修改前先检查 Git 状态，不覆盖或删除与当前任务无关的改动。
- 未经明确要求，不执行 commit、push、rebase、reset 或强制操作。

## 5. 完成标准

- `lint`、TypeScript 检查和生产构建通过。
- 手机与电脑布局可用，主要页面无明显溢出。
- 页面具备加载、空数据和错误状态，浏览器控制台无未处理错误。
