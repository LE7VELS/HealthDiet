# Project instructions

开始任何任务前先阅读 `PROJECT_CONTEXT.md`。

## 需要阅读的文档

- 前端任务：
  - `frontend/FRONTEND_REQUIREMENTS.md`
  - `frontend/FRONTEND_CONSTRAINTS.md`
  - 涉及接口时再读 `API_CONTRACT.md`
- 后端任务：
  - `backend/BACKEND_REQUIREMENTS.md`
  - `backend/BACKEND_CONSTRAINTS.md`
  - `API_CONTRACT.md`
  - `DATA_MODEL.md`
- 修改 API 或数据模型：同时检查前端和后端相关文档。

## 基本规则

- 保持小项目需要的简单结构，不提前增加复杂基础设施。
- 前后端代码都要包含充分且准确的中文注释。新增或修改模块时，至少检查包/模块职责、主要类型与 DTO、导出函数、关键业务规则、数据转换、错误映射、安全边界和不直观辅助函数是否已有说明；不能只给入口写一条概括性注释。
- 注释优先解释“为什么这样做、输入输出边界是什么、哪些值不能信任、错误如何转换”，不要求给简单赋值和显而易见的控制流逐行加注释。完成任务前逐个检查本次涉及文件，代码变化时同步维护注释，发现过时注释必须一并修正。
- 当前以手工功能测试为主；仍需保证格式化、静态检查和构建通过。
- 文档冲突时先指出，不自行猜测。
- 修改前检查 Git 状态，不覆盖无关改动。
- 不提交密钥、真实 `.env`、构建产物和本地上传文件。
- 未经明确要求，不执行 commit、push、rebase、reset 或强制操作。

## 本地 GitHub 工具

- 本机 `gh` 通过系统 Keyring 保存登录凭据；文档和日志不得记录账号、Token、仓库地址或凭据存储内容。
- Codex 受限环境可能因无法正确访问系统 Keyring，使 `gh auth status` 错误显示 `401 Bad credentials`。遇到该结果时，先在宿主权限下重新执行只读的 `gh auth status`，不能仅据此要求用户重新登录。
- 普通提交和推送直接使用 `git`；需要 PR、Issue 或 Actions 功能时，再按需使用宿主权限运行 `gh`。
