package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/LE7VELS/HealthDiet/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// currentUserIDKey 是 Handler 读取认证身份的唯一上下文键，值只能由 RequireAuth 写入。
const currentUserIDKey = "currentUserID"

// RequireAuth 验证 Bearer JWT，并只把签名保护的用户 ID 写入请求上下文。
func RequireAuth(tokens *service.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只接受标准的两段式 Authorization Header，拒绝缺失 Token、额外字段和其他认证方案。
		authorization := c.GetHeader("Authorization")
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			logAuthenticationRejection(c, "missing_or_invalid_authorization")
			abortUnauthenticated(c)
			return
		}
		// 用户 ID 只能来自完成签名与有效期校验的 JWT subject，不能信任请求参数中的 userId。
		userID, err := tokens.Verify(parts[1])
		if err != nil {
			logAuthenticationRejection(c, "invalid_token")
			abortUnauthenticated(c)
			return
		}
		c.Set(currentUserIDKey, userID)
		c.Next()
	}
}

// logAuthenticationRejection 只记录拒绝类别和请求元数据，不记录可能包含敏感信息的 Authorization 内容。
func logAuthenticationRejection(c *gin.Context, reason string) {
	log.Printf("JWT 认证失败 method=%s path=%s ip=%s code=UNAUTHENTICATED message=%q reason=%s", c.Request.Method, c.Request.URL.Path, c.ClientIP(), "登录状态无效或已过期", reason)
}

// CurrentUserID 从认证上下文读取当前用户 ID；返回 false 表示路由未经过认证或上下文被错误使用。
func CurrentUserID(c *gin.Context) (string, bool) {
	value, ok := c.Get(currentUserIDKey)
	userID, typeOK := value.(string)
	return userID, ok && typeOK && userID != ""
}

// abortUnauthenticated 统一终止请求，避免不同受保护接口泄露 Token 失败的具体原因。
func abortUnauthenticated(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{
		"code": "UNAUTHENTICATED", "message": "登录状态无效或已过期",
	}})
}
