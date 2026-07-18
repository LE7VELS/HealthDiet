# 前端开发约束

> 页面要求见 [`FRONTEND_REQUIREMENTS.md`](./FRONTEND_REQUIREMENTS.md)，接口见 [`../API_CONTRACT.md`](../API_CONTRACT.md)。

## 1. 技术栈

- React + TypeScript + Vite。
- Astryx Design System + StyleX。
- React Router。
- TanStack Query。
- React Hook Form + Zod。
- Axios。
- ECharts。

不增加 Redux 或第二套 UI、请求缓存方案，除非出现明确需求。

## 2. 代码结构

- 页面和路由放在 `app`、`pages`。
- 业务按 `features` 分为 auth、profile、foods、meals、nutrition。
- Astryx 组件通过 `components/ui` 封装。
- API 调用集中在 `lib/api`，页面不直接导入 Axios 或拼接 URL；`lib/api/client.ts` 统一配置基础地址、Bearer Token、响应解析和错误转换。
- 使用 TypeScript strict，避免 `any`。
- API DTO 和页面展示模型可以分开，但不要过度抽象。
- 食品和菜谱共用 `foods` 业务模块、页面与 API DTO，通过 Food 的 `kind` 区分，不复制两套状态和请求逻辑。
- 菜谱表单只提交原料 Food ID、用量、份数、成品重量和步骤；总营养、每份营养及每 100g 营养只展示后端结果。

### 中文注释

- 前端代码需要包含适量中文注释，说明组件职责、关键交互、状态流转、DTO 转换和不直观的兼容处理。
- 自定义 Hook、复杂表单联动、缓存刷新关系和容易误用的公共组件应说明设计意图或使用限制。
- 不给简单 JSX、明显的变量赋值和一眼可懂的样式逐行添加注释，避免注释淹没代码。
- 修改逻辑时同步更新相关注释，不保留与实际行为不一致的旧说明。

## 3. API 和状态

- 地址来自 `VITE_API_BASE_URL`。
- 业务 API 模块通过统一的 `apiRequest` 传递 Axios 请求配置；普通 JSON 请求使用 `data` 传入对象，由 Axios 负责序列化，不重复手写 `JSON.stringify`。
- 统一请求层允许业务接口补充必要 Header，但最终 `Authorization` 必须由会话层写入，业务调用方不得覆盖当前用户身份。
- 后端接口尚未完成时允许默认启用 Mock；接口实现并联调通过后，按接口或业务模块移除对应 Mock 分支，不要求在后端尚未就绪时全局关闭 Mock。
- Mock 和真实请求保持相同 DTO，页面组件不直接判断数据来源；最终主流程必须使用真实 Go API。
- 服务端数据使用 TanStack Query。
- Mutation 后刷新受影响的数据；饮食记录变化要刷新列表、Dashboard、每日汇总和趋势。
- 以后端营养计算和规则提示为准。
- 正确区分 `null` 和 `0`。
- 认证失败统一清理会话并跳转登录。
- 本地存在 Token 只表示浏览器保存了会话候选；进入受保护路由时必须通过 `/auth/me` 验证签名对应的服务端身份。验证成功可刷新本地公开用户信息，但不能重置 Token 的原始过期时间。

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
