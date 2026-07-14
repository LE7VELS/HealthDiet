package router

import (
	"fmt"

	"github.com/LE7VELS/HealthDiet/backend/internal/handler"
	"github.com/LE7VELS/HealthDiet/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Dependencies 是 Router 可以调用的 Handler 集合。
type Dependencies struct {
	RootHandler *handler.RootHandler
}

// New 创建 Gin Router，并按 Middleware → Handler 的顺序注册请求链。
func New(deps Dependencies) (*gin.Engine, error) {
	if deps.RootHandler == nil {
		return nil, fmt.Errorf("RootHandler 不能为空")
	}

	engine := gin.New()
	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	engine.Use(middleware.Logger(), gin.Recovery())
	engine.GET("/", deps.RootHandler.Get)
	engine.NoRoute(handler.NotFound)

	return engine, nil
}
