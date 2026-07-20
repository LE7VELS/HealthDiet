// Package apperr 集中定义需要跨 Store、Service 和 Handler 识别的稳定应用错误。
// 本包只表达错误身份，不依赖 Gin、HTTP 或 MongoDB Driver；调用方通过 errors.Is 判断，
// 需要补充内部定位上下文时必须使用 fmt.Errorf 的 %w 保留错误链。
package apperr

import "errors"

var (
	// ErrUserNotFound 表示按当前查询条件找不到用户；上层可按接口语义转换为登录失败、会话失效或资源不存在。
	ErrUserNotFound = errors.New("user not found")
	// ErrUserDuplicate 表示用户唯一字段触发持久化冲突，Service 必须继续转换为明确的用户名或邮箱冲突。
	ErrUserDuplicate = errors.New("user unique field duplicate")
	// ErrUsernameConflict 表示规范化后的用户名已被占用。
	ErrUsernameConflict = errors.New("username conflict")
	// ErrEmailConflict 表示规范化后的邮箱已被注册。
	ErrEmailConflict = errors.New("email conflict")
	// ErrInvalidCredentials 统一覆盖账号不存在和密码不匹配，避免泄露账号状态。
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidInput 标识请求字段校验失败；具体字段明细可由专用错误类型通过 Unwrap 关联到本错误。
	ErrInvalidInput = errors.New("invalid input")
	// ErrUnauthenticated 表示请求没有有效登录身份，包括无效 Token 和 Token 对应用户已不存在。
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrResourceNotFound 表示当前用户可见范围内的资源不存在。
	ErrResourceNotFound = errors.New("resource not found")
)
