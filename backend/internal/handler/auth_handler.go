package handler

// 本文件实现认证模块的 HTTP DTO、请求解析和成功响应转换；错误响应复用 error.go 的统一映射。

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/LE7VELS/HealthDiet/backend/internal/apperr"
	"github.com/LE7VELS/HealthDiet/backend/internal/middleware"
	"github.com/LE7VELS/HealthDiet/backend/internal/model"
	"github.com/LE7VELS/HealthDiet/backend/internal/service"
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

// NewAuthHandler 注入应用共享的认证 Service，Handler 自身不保存请求状态。
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

// Register 解析注册输入并返回 201；用户名、邮箱冲突和字段错误由统一错误结构表达。
func (h *AuthHandler) Register(c *gin.Context) {
	var request registerRequest
	if err := decodeJSON(c, &request); err != nil {
		writeAppError(c, apperr.ErrInvalidInput)
		return
	}

	// HTTP DTO 在边界处转换为 Service 输入，避免业务层依赖 Gin 或 JSON 标签。
	result, err := h.service.Register(c.Request.Context(), service.RegisterInput{
		Username: request.Username, Email: request.Email, Password: request.Password,
	})
	if err != nil {
		writeAppError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": authResultResponse(result)})
}

// Login 使用统一认证结果返回 JWT；失败时不区分账号不存在与密码错误，避免泄露账号状态。
func (h *AuthHandler) Login(c *gin.Context) {
	var request loginRequest
	if err := decodeJSON(c, &request); err != nil {
		writeAppError(c, apperr.ErrInvalidInput)
		return
	}

	result, err := h.service.Login(c.Request.Context(), request.Identifier, request.Password)
	if err != nil {
		writeAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": authResultResponse(result)})
}

// Me 根据 JWT 中间件建立的用户上下文查询当前账号，不能从查询参数接收 userId。
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		writeAppError(c, apperr.ErrUnauthenticated)
		return
	}
	user, err := h.service.CurrentUser(c.Request.Context(), userID)
	if err != nil {
		writeAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publicUserFromModel(user)})
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
