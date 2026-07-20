package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/LE7VELS/HealthDiet/backend/internal/apperr"
	"github.com/LE7VELS/HealthDiet/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// fieldError 描述某个请求字段的业务校验错误，字段名与 API 的 camelCase 请求字段保持一致。
type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// writeAppError 是应用错误到 HTTP 状态码、稳定 API 错误码和安全消息的统一映射点。
// 未识别错误只能返回通用内部错误，真实包装错误链仅写入服务端日志。
func writeAppError(c *gin.Context, err error) {
	var inputError *service.InputError
	switch {
	case errors.As(err, &inputError):
		fields := make([]fieldError, 0, len(inputError.Fields))
		for _, item := range inputError.Fields {
			fields = append(fields, fieldError{Field: item.Field, Message: item.Message})
		}
		logRequestFailure(c, "VALIDATION_ERROR", "请求参数不正确", nil)
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确", fields)
	case errors.Is(err, apperr.ErrInvalidInput):
		logRequestFailure(c, "VALIDATION_ERROR", "请求参数不正确", nil)
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确", nil)
	case errors.Is(err, apperr.ErrUsernameConflict):
		logRequestFailure(c, "USERNAME_CONFLICT", "用户名已被使用", nil)
		writeError(c, http.StatusConflict, "USERNAME_CONFLICT", "用户名已被使用", []fieldError{{Field: "username", Message: "用户名已被使用"}})
	case errors.Is(err, apperr.ErrEmailConflict):
		logRequestFailure(c, "EMAIL_CONFLICT", "邮箱已被注册", nil)
		writeError(c, http.StatusConflict, "EMAIL_CONFLICT", "邮箱已被注册", []fieldError{{Field: "email", Message: "邮箱已被注册"}})
	case errors.Is(err, apperr.ErrInvalidCredentials):
		logRequestFailure(c, "INVALID_CREDENTIALS", "账号或密码不正确", nil)
		writeError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "账号或密码不正确", nil)
	case errors.Is(err, apperr.ErrUnauthenticated):
		logRequestFailure(c, "UNAUTHENTICATED", "登录状态无效或已过期", nil)
		writeError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "登录状态无效或已过期", nil)
	case errors.Is(err, apperr.ErrUserNotFound), errors.Is(err, apperr.ErrResourceNotFound):
		logRequestFailure(c, "RESOURCE_NOT_FOUND", "请求的资源不存在", nil)
		writeError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "请求的资源不存在", nil)
	default:
		logRequestFailure(c, "INTERNAL_ERROR", "服务器暂时无法处理请求", err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器暂时无法处理请求", nil)
	}
}

// logRequestFailure 统一记录失败请求的安全元数据；只有未知内部错误附带完整包装错误链。
// 日志刻意不记录请求体、Authorization Header、密码或数据库连接信息。
func logRequestFailure(c *gin.Context, code, message string, err error) {
	if err != nil {
		log.Printf("请求失败 method=%s path=%s ip=%s code=%s message=%q error=%v", c.Request.Method, c.Request.URL.Path, c.ClientIP(), code, message, err)
		return
	}
	log.Printf("请求失败 method=%s path=%s ip=%s code=%s message=%q", c.Request.Method, c.Request.URL.Path, c.ClientIP(), code, message)
}

// writeError 生成 API_CONTRACT.md 规定的统一错误包络；无字段错误时省略 fields。
func writeError(c *gin.Context, status int, code, message string, fields []fieldError) {
	errorBody := gin.H{"code": code, "message": message}
	if len(fields) > 0 {
		errorBody["fields"] = fields
	}
	c.JSON(status, gin.H{"error": errorBody})
}
