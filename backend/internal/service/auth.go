package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/LE7VELS/HealthDiet/backend/internal/model"
	"github.com/LE7VELS/HealthDiet/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUsernameConflict 表示规范化后的用户名已被占用。
	ErrUsernameConflict = errors.New("用户名已被使用")
	// ErrEmailConflict 表示规范化后的邮箱已被注册。
	ErrEmailConflict = errors.New("邮箱已被注册")
	// ErrInvalidCredentials 统一覆盖账号不存在和密码不匹配，避免账号枚举。
	ErrInvalidCredentials = errors.New("账号或密码不正确")
	// ErrInvalidInput 是所有字段校验错误的通用标识，详细字段保存在 InputError 中。
	ErrInvalidInput = errors.New("认证参数不正确")
)

// RegisterInput 是注册业务所需的最小输入，不包含 HTTP、JSON 或 MongoDB 实现细节。
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// InputField 是业务层字段校验结果；字段名使用 API camelCase，供 Handler 原样映射。
type InputField struct {
	Field   string
	Message string
}

// InputError 保留字段级校验结果，供 Handler 转换为统一 API 错误结构。
type InputError struct {
	Fields []InputField
}

// Error 实现 error 接口，并返回不包含用户输入和内部细节的稳定消息。
func (e *InputError) Error() string { return ErrInvalidInput.Error() }

// Unwrap 允许调用方通过 errors.Is 识别通用输入错误，同时仍可通过 errors.As 取得字段明细。
func (e *InputError) Unwrap() error { return ErrInvalidInput }

// AuthResult 汇总签发后的访问 Token 与用户模型；Handler 必须裁剪 User 后再返回客户端。
type AuthResult struct {
	AccessToken string
	ExpiresIn   int64
	User        model.User
}

// AuthService 编排冲突检查、密码哈希和 Token 签发，不向 HTTP 层暴露密码哈希。
type AuthService struct {
	store  *store.Store
	tokens *TokenManager
}

// NewAuthService 注入认证流程所需的唯一 Store 和 TokenManager。
func NewAuthService(appStore *store.Store, tokens *TokenManager) *AuthService {
	return &AuthService{store: appStore, tokens: tokens}
}

// Register 完成输入规范化、业务校验、唯一性检查、密码哈希、持久化和 Token 签发。
// 明文密码只在当前调用栈中短暂使用，进入 Store 前必须转换为 bcrypt 哈希。
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	// 用户名保留展示大小写，邮箱统一小写保存；唯一查询还会在 Store 中再次规范化。
	username := strings.TrimSpace(input.Username)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	fields := make([]InputField, 0, 3)
	if !validUsername(username) {
		fields = append(fields, InputField{Field: "username", Message: "用户名需为 2 至 24 个文字、数字、下划线或连字符"})
	}
	if !validEmail(email) {
		fields = append(fields, InputField{Field: "email", Message: "邮箱格式不正确"})
	}
	if !validPassword(input.Password) {
		fields = append(fields, InputField{Field: "password", Message: "密码需为 8 至 64 个字符，并包含字母和数字"})
	}
	if len(fields) > 0 {
		return AuthResult{}, &InputError{Fields: fields}
	}

	// 提前检查用于返回明确的冲突字段；MongoDB 唯一索引仍是并发写入时的最终保证。
	if _, err := s.store.FindUserByUsername(ctx, username); err == nil {
		return AuthResult{}, ErrUsernameConflict
	} else if !errors.Is(err, store.ErrUserNotFound) {
		return AuthResult{}, err
	}
	if _, err := s.store.FindUserByEmail(ctx, email); err == nil {
		return AuthResult{}, ErrEmailConflict
	} else if !errors.Is(err, store.ErrUserNotFound) {
		return AuthResult{}, err
	}

	// bcrypt 自带盐并使用适合密码的计算成本，数据库永远不接收明文密码。
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, fmt.Errorf("生成密码哈希: %w", err)
	}
	user, err := s.store.CreateUser(ctx, username, email, string(hash))
	if err != nil {
		// 唯一索引是并发注册的最终安全边界；重复时再次查询以返回稳定的业务错误。
		if errors.Is(err, store.ErrUserDuplicate) {
			if _, findErr := s.store.FindUserByUsername(ctx, username); findErr == nil {
				return AuthResult{}, ErrUsernameConflict
			}
			return AuthResult{}, ErrEmailConflict
		}
		return AuthResult{}, err
	}
	return s.resultFor(user)
}

// Login 使用统一的无效凭证错误处理“不存在用户”和“密码不匹配”，避免账号枚举。
func (s *AuthService) Login(ctx context.Context, identifier, password string) (AuthResult, error) {
	if strings.TrimSpace(identifier) == "" || password == "" {
		return AuthResult{}, ErrInvalidCredentials
	}
	user, err := s.store.FindUserByIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, err
	}
	// bcrypt 比较负责处理哈希格式和恒定成本，不能自行比较字符串或重新哈希明文。
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return AuthResult{}, ErrInvalidCredentials
	}
	return s.resultFor(user)
}

// CurrentUser 通过认证后的不透明用户 ID 读取账号，供 /auth/me 使用。
func (s *AuthService) CurrentUser(ctx context.Context, userID string) (model.User, error) {
	return s.store.FindUserByID(ctx, userID)
}

// resultFor 集中签发登录与注册共用的访问 Token，避免两条流程产生不同会话语义。
func (s *AuthService) resultFor(user model.User) (AuthResult, error) {
	token, err := s.tokens.Create(user.ID)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{AccessToken: token, ExpiresIn: s.tokens.ExpiresInSeconds(), User: user}, nil
}

// validUsername 按 Unicode 字符数校验用户名，允许中英文、数字、下划线和连字符。
func validUsername(value string) bool {
	length := utf8.RuneCountInString(value)
	if length < 2 || length > 24 {
		return false
	}
	for _, character := range value {
		if !unicode.IsLetter(character) && !unicode.IsNumber(character) && character != '_' && character != '-' {
			return false
		}
	}
	return true
}

// validEmail 要求输入是单一邮箱地址而不是带显示名称的邮件地址，并限制标准最大长度。
func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && len(value) <= 254
}

// validPassword 执行服务端基础复杂度校验；客户端校验只能改善体验，不能替代此边界。
func validPassword(value string) bool {
	if len(value) < 8 || len(value) > 64 {
		return false
	}
	hasLetter, hasDigit := false, false
	for _, character := range value {
		hasLetter = hasLetter || unicode.IsLetter(character)
		hasDigit = hasDigit || unicode.IsDigit(character)
	}
	return hasLetter && hasDigit
}
