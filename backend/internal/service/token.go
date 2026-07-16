package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// tokenIssuer 防止其他系统使用同一密钥签发的 Token 被误接受为本 API 的凭证。
const tokenIssuer = "healthdiet-api"

// ErrInvalidToken 对外统一表示格式、算法、签名或时间声明校验失败，不泄露具体失败步骤。
var ErrInvalidToken = errors.New("Token 无效或已过期")

// tokenClaims 是当前访问 Token 的最小声明集；subject 只保存用户 ID，不放入密码或隐私资料。
type tokenClaims struct {
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// TokenManager 使用 HS256 签发和验证短期访问 Token。
// 密钥只保存在进程内存中，Token 本身不持久化到 users 集合。
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenManager 创建可复用的 Token 管理器；密钥长度在 Config.Load 阶段统一校验。
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

// ExpiresInSeconds 返回 API 合同使用的秒数，避免 Handler 了解内部 time.Duration。
func (m *TokenManager) ExpiresInSeconds() int64 {
	return int64(m.ttl / time.Second)
}

// Create 为指定用户签发 JWT，并写入签发时间、过期时间和固定 issuer。
func (m *TokenManager) Create(userID string) (string, error) {
	now := time.Now().UTC()
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", fmt.Errorf("序列化 JWT Header: %w", err)
	}
	claims, err := json.Marshal(tokenClaims{
		Subject: userID, Issuer: tokenIssuer, IssuedAt: now.Unix(), ExpiresAt: now.Add(m.ttl).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("序列化 JWT Claims: %w", err)
	}

	// JWT 签名覆盖编码后的 Header 与 Claims，任一部分被修改都会导致验签失败。
	unsigned := encodeTokenPart(header) + "." + encodeTokenPart(claims)
	return unsigned + "." + encodeTokenPart(m.sign(unsigned)), nil
}

// Verify 严格校验 JWT 三段格式、固定算法、HMAC 签名、issuer、subject 和时间窗口。
// 只有全部通过后才返回可写入认证上下文的用户 ID。
func (m *TokenManager) Verify(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidToken
	}

	// 先固定算法为 HS256，避免攻击者通过 Header 要求服务端使用其他或无签名算法。
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalidToken
	}
	var header map[string]string
	if json.Unmarshal(headerBytes, &header) != nil || header["alg"] != "HS256" || header["typ"] != "JWT" {
		return "", ErrInvalidToken
	}

	// hmac.Equal 使用恒定时间比较，避免普通字节比较泄露签名匹配进度。
	providedSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(providedSignature, m.sign(parts[0]+"."+parts[1])) {
		return "", ErrInvalidToken
	}

	// 签名通过后再信任 Claims；subject 为空或 issuer 不匹配都不能建立用户身份。
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidToken
	}
	var claims tokenClaims
	if json.Unmarshal(claimsBytes, &claims) != nil || claims.Issuer != tokenIssuer || claims.Subject == "" {
		return "", ErrInvalidToken
	}
	// 允许最多 60 秒服务器时钟偏差，但不接受已过期或明显来自未来的 Token。
	now := time.Now().Unix()
	if claims.ExpiresAt <= now || claims.IssuedAt > now+60 {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}

// sign 计算 HS256 的 HMAC-SHA256 摘要；调用方负责对结果做 JWT Base64URL 编码。
func (m *TokenManager) sign(value string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(value))
	return mac.Sum(nil)
}

// encodeTokenPart 使用 JWT 要求的无填充 Base64URL 编码，不能替换为普通 Base64。
func encodeTokenPart(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}
