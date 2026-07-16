// Package main 负责装配并启动 HealthDiet 的单个 Gin API 进程。
// 业务规则留在 Service，数据库访问留在 Store，入口只管理配置、依赖和生命周期。
package main

import (
	"context"
	"log"
	"time"

	"github.com/LE7VELS/HealthDiet/backend/internal/config"
	"github.com/LE7VELS/HealthDiet/backend/internal/handler"
	"github.com/LE7VELS/HealthDiet/backend/internal/router"
	"github.com/LE7VELS/HealthDiet/backend/internal/service"
	"github.com/LE7VELS/HealthDiet/backend/internal/store"
)

const (
	// MongoDB 的启动与关闭都使用有限超时，避免外部依赖异常时进程永久阻塞。
	mongoStartupTimeout  = 10 * time.Second
	mongoShutdownTimeout = 5 * time.Second
	// 当前合同只提供短期访问 Token，不提供刷新 Token，因此固定一小时有效期。
	accessTokenTTL = time.Hour
)

// main 只负责把可报告的启动错误交给标准日志并终止进程。
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run 按“配置 → 数据库 → Service/Handler → Router”的依赖方向组装单个 API 应用。
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 先确认数据库和必要索引可用，再启动 HTTP 服务，避免接受无法完成的业务请求。
	startupContext, cancelStartup := context.WithTimeout(context.Background(), mongoStartupTimeout)
	appStore, err := store.New(startupContext, cfg.MongoDBURI, cfg.MongoDBName)
	if err != nil {
		cancelStartup()
		return err
	}

	err = appStore.EnsureSchema(startupContext)
	cancelStartup()
	if err != nil {
		closeStore(appStore)
		return err
	}

	// Store 由应用入口统一持有，后续业务 Handler 或 Service 共享同一个 MongoDB Client。
	defer closeStore(appStore)

	// TokenManager 与 AuthService 在应用生命周期内复用，避免 Handler 自行读取密钥或访问数据库。
	tokens := service.NewTokenManager(cfg.JWTSecret, accessTokenTTL)
	authService := service.NewAuthService(appStore, tokens)
	authHandler := handler.NewAuthHandler(authService)
	// Router 使用命名依赖字段装配各模块，后续增加 Handler 时不需要不断扩展 New 的位置参数。
	engine, err := router.New(router.Dependencies{
		AuthHandler:  authHandler,
		TokenManager: tokens,
		CORSOrigin:   cfg.CORSOrigin,
	})
	if err != nil {
		return err
	}

	log.Printf("MongoDB 已连接，数据库 %s", cfg.MongoDBName)
	log.Printf("Gin API 已启动，监听地址 %s", cfg.HTTPAddress)

	return engine.Run(cfg.HTTPAddress)
}

// closeStore 关闭应用共享的 MongoDB Client。
func closeStore(appStore *store.Store) {
	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), mongoShutdownTimeout)
	defer cancelShutdown()

	if err := appStore.Close(shutdownContext); err != nil {
		log.Printf("关闭 MongoDB 连接失败: %v", err)
	}
}
