package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppService 是 RootHandler 需要的 Service 能力。
type AppService interface {
	Message(context.Context) string
}

// RootHandler 处理根路径的 HTTP 输入输出。
type RootHandler struct {
	service AppService
}

func NewRootHandler(service AppService) *RootHandler {
	return &RootHandler{service: service}
}

func (h *RootHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": h.service.Message(c.Request.Context()),
	})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": gin.H{
			"code":    "RESOURCE_NOT_FOUND",
			"message": "请求的资源不存在",
		},
	})
}
