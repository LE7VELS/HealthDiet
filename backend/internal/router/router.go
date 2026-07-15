package router

import (
	"github.com/LE7VELS/HealthDiet/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

// New 创建 Gin Router，并按 Middleware → Handler 的顺序注册请求链。
func New() (*gin.Engine, error) {
	engine := gin.New()
	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	engine.Use(gin.Logger(), gin.Recovery())
	engine.GET("/", handler.Root)
	engine.NoRoute(handler.NotFound)

	return engine, nil
}
