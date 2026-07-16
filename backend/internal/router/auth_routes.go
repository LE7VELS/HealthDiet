package router

import (
	"github.com/LE7VELS/HealthDiet/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// registerAuthRoutes 注册认证模块路由。
// 注册和登录是公开入口；/me 必须先由 JWT 中间件建立可信用户上下文。
func registerAuthRoutes(api *gin.RouterGroup, deps Dependencies) {
	auth := api.Group("/auth")
	auth.POST("/register", deps.AuthHandler.Register)
	auth.POST("/login", deps.AuthHandler.Login)
	auth.GET("/me", middleware.RequireAuth(deps.TokenManager), deps.AuthHandler.Me)
}
