package router

import (
	"github.com/LE7VELS/HealthDiet/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

// registerRootRoutes 注册不属于版本化业务 API 的开发确认入口和统一 404 Handler。
func registerRootRoutes(engine *gin.Engine) {
	engine.GET("/", handler.Root)
	engine.NoRoute(handler.NotFound)
}
