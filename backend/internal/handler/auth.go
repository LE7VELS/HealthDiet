package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/LE7VELS/HealthDiet/backend/internal/middleware"
	"github.com/LE7VELS/HealthDiet/backend/internal/model"
	"github.com/LE7VELS/HealthDiet/backend/internal/service"
	"github.com/LE7VELS/HealthDiet/backend/internal/store"
	"github.com/gin-gonic/gin"
)

// AuthHandler 把认证 HTTP 请求交给 AuthService，并统一转换成功响应与业务错误。
type AuthHandler struct {
	service *service.AuthService
}

// registerRequest 是注册接口允许接收的完整字段集合；额外 JSON 字段会被 decodeJSON 拒绝。
type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginRequest 支持用户名或邮箱作为 identifier，具体规范化与密码校验由 Service 完成。
type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// publicUser 是可返回给浏览器的最小用户视图，刻意不包含密码哈希和数据库时间字段。
type publicUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// authResponse 与 API_CONTRACT.md 的登录/注册响应保持一致。
// ExpiresIn 使用秒，前端据此判断本地会话是否过期。
type authResponse struct {
	AccessToken string     `json:"accessToken"`
	TokenType   string     `json:"tokenType"`
	ExpiresIn   int64      `json:"expiresIn"`
	User        publicUser `json:"user"`
}

// fieldError 描述某个请求字段的业务校验错误，便于前端回填到对应输入框。
type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewAuthHandler 注入应用共享的认证 Service，Handler 自身不保存请求状态。
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

// Register 解析注册输入并返回 201；用户名、邮箱冲突和字段错误由统一错误结构表达。
func (h *AuthHandler) Register(c *gin.Context) {
	var request registerRequest
	if err := decodeJSON(c, &request); err != nil {
		logAuthFailure(c, "VALIDATION_ERROR", "请求参数不正确", err)
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确", nil)
		return
	}

	// HTTP DTO 在边界处转换为 Service 输入，避免业务层依赖 Gin 或 JSON 标签。
	result, err := h.service.Register(c.Request.Context(), service.RegisterInput{
		Username: request.Username, Email: request.Email, Password: request.Password,
	})
	if err != nil {
		h.writeAuthError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": authResultResponse(result)})
}

// Login 使用统一认证结果返回 JWT；失败时不区分账号不存在与密码错误，避免泄露账号状态。
func (h *AuthHandler) Login(c *gin.Context) {
	var request loginRequest
	if err := decodeJSON(c, &request); err != nil {
		logAuthFailure(c, "VALIDATION_ERROR", "请求参数不正确", err)
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确", nil)
		return
	}

	result, err := h.service.Login(c.Request.Context(), request.Identifier, request.Password)
	if err != nil {
		h.writeAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": authResultResponse(result)})
}

// Me 根据 JWT 中间件建立的用户上下文查询当前账号，不能从查询参数接收 userId。
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		logAuthFailure(c, "UNAUTHENTICATED", "请先登录", nil)
		writeError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "请先登录", nil)
		return
	}
	user, err := h.service.CurrentUser(c.Request.Context(), userID)
	if err != nil {
		// Token 指向的用户已不存在时视为会话失效；其他 Store 错误统一隐藏为内部错误。
		if errors.Is(err, store.ErrUserNotFound) {
			logAuthFailure(c, "UNAUTHENTICATED", "登录状态已失效", nil)
			writeError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "登录状态已失效", nil)
			return
		}
		logAuthFailure(c, "INTERNAL_ERROR", "服务器暂时无法处理请求", err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器暂时无法处理请求", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publicUserFromModel(user)})
}

// writeAuthError 是认证业务错误到 HTTP 状态码和稳定错误码的唯一映射点。
// 未识别的内部错误不得把 MongoDB、哈希或密钥细节返回给客户端。
func (h *AuthHandler) writeAuthError(c *gin.Context, err error) {
	var inputError *service.InputError
	switch {
	case errors.As(err, &inputError):
		logAuthFailure(c, "VALIDATION_ERROR", "请求参数不正确", nil)
		fields := make([]fieldError, 0, len(inputError.Fields))
		for _, item := range inputError.Fields {
			fields = append(fields, fieldError{Field: item.Field, Message: item.Message})
		}
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不正确", fields)
	case errors.Is(err, service.ErrUsernameConflict):
		logAuthFailure(c, "USERNAME_CONFLICT", "用户名已被使用", nil)
		writeError(c, http.StatusConflict, "USERNAME_CONFLICT", "用户名已被使用", []fieldError{{Field: "username", Message: "用户名已被使用"}})
	case errors.Is(err, service.ErrEmailConflict):
		logAuthFailure(c, "EMAIL_CONFLICT", "邮箱已被注册", nil)
		writeError(c, http.StatusConflict, "EMAIL_CONFLICT", "邮箱已被注册", []fieldError{{Field: "email", Message: "邮箱已被注册"}})
	case errors.Is(err, service.ErrInvalidCredentials):
		logAuthFailure(c, "INVALID_CREDENTIALS", "账号或密码不正确", nil)
		writeError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "账号或密码不正确", nil)
	default:
		logAuthFailure(c, "INTERNAL_ERROR", "服务器暂时无法处理请求", err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器暂时无法处理请求", nil)
	}
}

// logAuthFailure 记录认证请求的稳定错误码和人类可读消息。
// 内部错误额外记录完整包装链，但不记录请求体、密码或 JWT。
func logAuthFailure(c *gin.Context, code, message string, err error) {
	if err != nil {
		log.Printf("认证请求失败 method=%s path=%s ip=%s code=%s message=%q error=%v", c.Request.Method, c.Request.URL.Path, c.ClientIP(), code, message, err)
		return
	}
	log.Printf("认证请求失败 method=%s path=%s ip=%s code=%s message=%q", c.Request.Method, c.Request.URL.Path, c.ClientIP(), code, message)
}

// decodeJSON 只接受单个 JSON 对象并拒绝未知字段，防止客户端误以为服务端处理了未支持的参数。
func decodeJSON(c *gin.Context, target any) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("请求体只能包含一个 JSON 对象")
	}
	return nil
}

// authResultResponse 将包含内部 User Model 的 Service 结果裁剪为公开认证响应。
func authResultResponse(result service.AuthResult) authResponse {
	return authResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   result.ExpiresIn,
		User:        publicUserFromModel(result.User),
	}
}

// publicUserFromModel 集中维护用户模型到公开 DTO 的白名单字段转换。
func publicUserFromModel(user model.User) publicUser {
	return publicUser{ID: user.ID, Username: user.Username, Email: user.Email}
}

// writeError 生成 API_CONTRACT.md 规定的统一错误包络；无字段错误时省略 fields。
func writeError(c *gin.Context, status int, code, message string, fields []fieldError) {
	errorBody := gin.H{"code": code, "message": message}
	if len(fields) > 0 {
		errorBody["fields"] = fields
	}
	c.JSON(status, gin.H{"error": errorBody})
}
