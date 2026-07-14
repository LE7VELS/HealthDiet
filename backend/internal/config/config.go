package config

import "os"

const defaultHTTPAddress = ":8080"

// Config 保存应用配置；后续 MongoDB 和 JWT 配置也集中放在这里。
type Config struct {
	HTTPAddress string
}

// Load 从环境变量读取配置。
func Load() Config {
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = defaultHTTPAddress
	}

	return Config{HTTPAddress: address}
}
