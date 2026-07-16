// Package config 集中读取和校验进程启动配置。
// 本地 .env 只用于开发，部署环境变量具有更高优先级，敏感配置不提供代码默认值。
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	// 这些默认值仅用于本地开发；JWT 密钥刻意不提供默认值，防止部署时误用固定密钥。
	defaultHTTPAddress = ":8080"
	defaultMongoDBURI  = "mongodb://127.0.0.1:27017"
	defaultMongoDBName = "healthdiet"
	defaultCORSOrigin  = "http://localhost:5173"
)

// Config 保存启动后不再变化的应用配置，敏感值只从运行环境或本地 .env 读取。
type Config struct {
	// HTTPAddress 是 Gin 监听地址，例如 :8080。
	HTTPAddress string
	// MongoDBURI 和 MongoDBName 共同确定业务数据库，连接串不得返回给客户端。
	MongoDBURI  string
	MongoDBName string
	// JWTSecret 只用于服务端签名和验签，禁止写入日志或 API 响应。
	JWTSecret string
	// CORSOrigin 是唯一允许直接跨域访问 API 的前端来源。
	CORSOrigin string
}

// Load 先读取本地 .env，再读取配置；进程中已有的环境变量不会被 .env 覆盖。
func Load() (Config, error) {
	// godotenv.Load 不覆盖进程中已有变量，因此部署平台注入的值天然优先于本地文件。
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("读取 .env: %w", err)
	}

	config := Config{
		HTTPAddress: valueOrDefault("HTTP_ADDR", defaultHTTPAddress),
		MongoDBURI:  valueOrDefault("MONGODB_URI", defaultMongoDBURI),
		MongoDBName: valueOrDefault("MONGODB_DATABASE", defaultMongoDBName),
		JWTSecret:   strings.TrimSpace(os.Getenv("JWT_SECRET")),
		CORSOrigin:  valueOrDefault("CORS_ALLOWED_ORIGIN", defaultCORSOrigin),
	}
	if len(config.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET 必须通过环境变量设置，且至少包含 32 个字符")
	}

	return config, nil
}

// valueOrDefault 只为非敏感开发配置提供回退值，JWT_SECRET 不得通过此函数设置默认值。
func valueOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		// 默认值只服务本地开发，部署环境通过环境变量覆盖连接信息。
		return fallback
	}

	return value
}
