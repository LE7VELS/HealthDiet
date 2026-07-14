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
- 当前以手工功能测试为主；仍需保证格式化、静态检查和构建通过。
- 文档冲突时先指出，不自行猜测。
- 修改前检查 Git 状态，不覆盖无关改动。
- 不提交密钥、真实 `.env`、构建产物和本地上传文件。
- 未经明确要求，不执行 commit、push、rebase、reset 或强制操作。
