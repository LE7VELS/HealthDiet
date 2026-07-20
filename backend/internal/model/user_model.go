package model

// 本文件定义认证和用户查询流程共享的用户业务模型。

import "time"

// User 是认证流程在 Handler、Service 和 Store 之间共享的业务模型。
// ID 使用不透明字符串隔离 MongoDB ObjectID；PasswordHash 仅供认证 Service 校验，Handler 必须转换为公开 DTO 后才能响应。
type User struct {
	// ID、Username 和 Email 是账号身份字段。
	ID       string
	Username string
	Email    string
	// PasswordHash 保存 bcrypt 结果，任何日志或 API 响应都不得包含此字段。
	PasswordHash string
	// CreatedAt 和 UpdatedAt 统一使用 UTC 时间。
	CreatedAt time.Time
	UpdatedAt time.Time
}
