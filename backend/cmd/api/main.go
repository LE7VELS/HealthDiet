package main

import (
	"log"

	"github.com/LE7VELS/HealthDiet/backend/internal/config"
	"github.com/LE7VELS/HealthDiet/backend/internal/handler"
	"github.com/LE7VELS/HealthDiet/backend/internal/router"
	"github.com/LE7VELS/HealthDiet/backend/internal/service"
)

func main() {
	cfg := config.Load()

	appService := service.NewAppService("HealthDiet API")
	rootHandler := handler.NewRootHandler(appService)

	engine, err := router.New(router.Dependencies{
		RootHandler: rootHandler,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Gin API 已启动，监听地址 %s", cfg.HTTPAddress)
	if err := engine.Run(cfg.HTTPAddress); err != nil {
		log.Fatal(err)
	}
}
