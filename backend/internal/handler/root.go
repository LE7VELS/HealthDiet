package handler

import (
	"net/http"

	"github.com/LE7VELS/HealthDiet/backend/internal/apperr"
	"github.com/gin-gonic/gin"
)

// Root 用于开发时手工确认 Gin 已启动。
func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "HealthDiet API",
	})
}

// NotFound 为所有未注册路径返回统一资源不存在错误，避免 Gin 默认的纯文本 404。
func NotFound(c *gin.Context) {
	writeAppError(c, apperr.ErrResourceNotFound)
}
