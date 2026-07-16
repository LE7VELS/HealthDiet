// Package router 负责注册版本化路由及其 Middleware → Handler 请求链。
// 本包不承载业务规则，也不直接调用 Service 或 Store。
package router

import (
	"fmt"

	"github.com/LE7VELS/HealthDiet/backend/internal/handler"
	"github.com/LE7VELS/HealthDiet/backend/internal/middleware"
	"github.com/LE7VELS/HealthDiet/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// Dependencies 汇总 Router 注册各业务路由所需的共享依赖。
// 使用命名字段传入可以避免 New 的参数列表随业务模块增长，并让 main.go 的装配关系保持可读。
type Dependencies struct {
	// AuthHandler 处理公开注册、登录和当前用户查询。
	AuthHandler *handler.AuthHandler
	// TokenManager 提供受保护路由共用的 JWT 验证能力。
	TokenManager *service.TokenManager
	// CORSOrigin 是浏览器跨域访问时唯一允许的前端来源。
	CORSOrigin string
}

// New 创建 Gin Router，并集中注册全局中间件和各模块路由。
// Router 只做 HTTP 链路装配，不直接访问 Store 或承载认证业务规则。
func New(deps Dependencies) (*gin.Engine, error) {
	if err := validateDependencies(deps); err != nil {
		return nil, err
	}

	engine := gin.New()
	// 当前不信任任何反向代理，避免客户端伪造代理头影响来源地址判断。
	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	// 日志、异常恢复和 CORS 作用于全部路由；JWT 只挂载到需要登录的具体路由或路由组。
	engine.Use(gin.Logger(), gin.Recovery(), middleware.CORS(deps.CORSOrigin))
	registerRootRoutes(engine)

	// 所有业务接口统一放在 /api/v1，避免前端依赖没有版本的临时路径。
	api := engine.Group("/api/v1")
	registerAuthRoutes(api, deps)

	return engine, nil
}

// validateDependencies 在服务启动时发现装配错误，避免缺失依赖直到首个请求才触发空指针异常。
func validateDependencies(deps Dependencies) error {
	if deps.AuthHandler == nil {
		return fmt.Errorf("router: AuthHandler 不能为空")
	}
	if deps.TokenManager == nil {
		return fmt.Errorf("router: TokenManager 不能为空")
	}
	if deps.CORSOrigin == "" {
		return fmt.Errorf("router: CORSOrigin 不能为空")
	}
	return nil
}
