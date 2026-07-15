package main

import (
	"log"

	"github.com/LE7VELS/HealthDiet/backend/internal/config"
	"github.com/LE7VELS/HealthDiet/backend/internal/router"
)

func main() {
	cfg := config.Load()

	engine, err := router.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Gin API 已启动，监听地址 %s", cfg.HTTPAddress)
	if err := engine.Run(cfg.HTTPAddress); err != nil {
		log.Fatal(err)
	}
}
