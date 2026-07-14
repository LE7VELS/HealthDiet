# 前端开发约束

> 页面要求见 [`FRONTEND_REQUIREMENTS.md`](./FRONTEND_REQUIREMENTS.md)，接口见 [`../API_CONTRACT.md`](../API_CONTRACT.md)。

## 1. 技术栈

- React + TypeScript + Vite。
- Astryx Design System + StyleX。
- React Router。
- TanStack Query。
- React Hook Form + Zod。
- ECharts。

不增加 Redux 或第二套 UI、请求缓存方案，除非出现明确需求。

## 2. 代码结构

- 页面和路由放在 `app`、`pages`。
- 业务按 `features` 分为 auth、profile、foods、meals、nutrition。
- Astryx 组件通过 `components/ui` 封装。
- API 调用集中在 `lib/api`，页面不直接拼 URL。
- 使用 TypeScript strict，避免 `any`。
- API DTO 和页面展示模型可以分开，但不要过度抽象。

## 3. API 和状态

- 地址来自 `VITE_API_BASE_URL`。
- 服务端数据使用 TanStack Query。
- Mutation 后刷新受影响的数据；饮食记录变化要刷新列表、Dashboard、每日汇总和趋势。
- 以后端营养计算和规则提示为准。
- 正确区分 `null` 和 `0`。
- 认证失败统一清理会话并跳转登录。

## 4. 表单和上传

- 表单使用 React Hook Form + Zod。
- 客户端校验用于改善体验，服务端校验才是最终结果。
- 提交中禁止重复提交，失败后保留用户输入。
- 图片支持电脑选择或拖拽、手机拍照或相册选择。
- 每条记录最多一张图片，只允许 JPEG、PNG、WebP，最大 10 MB。
- 前端不得保存 OSS 密钥或服务端文件路径。

## 5. 样式和交互

- 同时适配手机和电脑，不做单纯缩放。
- 所有数据页有加载、空数据和错误状态。
- 成功使用 Toast，删除使用确认对话框。
- 表单有标签，主要操作可用键盘完成。

## 6. 检查方式

当前以手工功能测试为主，重点检查主流程、错误提示、图片上传和响应式布局。

提交前执行：

- TypeScript 检查。
- lint。
- 生产构建。

## 7. Git

- 不提交 `node_modules`、`dist`、真实 `.env`、密钥、日志和缓存。
- 不覆盖无关改动。
- 未经授权不执行 commit、push、rebase、reset 或强制操作。
