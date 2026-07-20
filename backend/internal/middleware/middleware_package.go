// Package middleware 保存跨路由共享的 HTTP 边界逻辑，目前包括 JWT 认证和 CORS。
// 中间件只验证请求并建立可信上下文，不查询数据库、不承载注册登录等业务规则；通用日志使用 Gin 官方 gin.Logger。
package middleware
