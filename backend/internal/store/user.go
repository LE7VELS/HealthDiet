package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/LE7VELS/HealthDiet/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	// ErrUserNotFound 隐藏 MongoDB 的 ErrNoDocuments，避免上层依赖 Driver 错误类型。
	ErrUserNotFound = errors.New("用户不存在")
	// ErrUserDuplicate 表示用户名或邮箱触发唯一索引，具体冲突字段由 Service 判定。
	ErrUserDuplicate = errors.New("用户唯一字段重复")
)

// userDocument 精确描述 users 集合的 BSON 结构，只在 Store 边界内使用。
// 规范化字段服务于不区分大小写的唯一查询，PasswordHash 永远不能转换到公开 DTO。
type userDocument struct {
	ID                 bson.ObjectID `bson:"_id"`
	Username           string        `bson:"username"`
	UsernameNormalized string        `bson:"username_normalized"`
	Email              string        `bson:"email"`
	EmailNormalized    string        `bson:"email_normalized"`
	PasswordHash       string        `bson:"password_hash"`
	CreatedAt          time.Time     `bson:"created_at"`
	UpdatedAt          time.Time     `bson:"updated_at"`
}

// CreateUser 保存已经完成业务校验和密码哈希的用户，并由 Store 生成 ObjectID 与 UTC 时间。
// 调用方不得传入明文密码；唯一索引冲突会转换为 ErrUserDuplicate。
func (s *Store) CreateUser(ctx context.Context, username, email, passwordHash string) (model.User, error) {
	now := time.Now().UTC()
	document := userDocument{
		ID:                 bson.NewObjectID(),
		Username:           username,
		UsernameNormalized: normalizeIdentity(username),
		Email:              email,
		EmailNormalized:    normalizeIdentity(email),
		PasswordHash:       passwordHash,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if _, err := s.database.Collection("users").InsertOne(ctx, document); err != nil {
		// Driver 的重复键错误在数据访问边界转换，上层不需要导入 mongo 包。
		if mongo.IsDuplicateKeyError(err) {
			return model.User{}, ErrUserDuplicate
		}
		return model.User{}, fmt.Errorf("创建用户: %w", err)
	}

	return userDocumentToModel(document), nil
}

// FindUserByUsername 使用规范化唯一字段查询，供注册冲突检查使用。
func (s *Store) FindUserByUsername(ctx context.Context, username string) (model.User, error) {
	return s.findUser(ctx, bson.D{{Key: "username_normalized", Value: normalizeIdentity(username)}})
}

// FindUserByEmail 使用规范化唯一字段查询，邮箱大小写差异不会绕过唯一约束。
func (s *Store) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	return s.findUser(ctx, bson.D{{Key: "email_normalized", Value: normalizeIdentity(email)}})
}

// FindUserByIdentifier 支持登录时使用用户名或邮箱，但不接受其他模糊匹配。
func (s *Store) FindUserByIdentifier(ctx context.Context, identifier string) (model.User, error) {
	// 用户名和邮箱共享同一规范化方式，因此登录只需要一次精确的 $or 查询。
	normalized := normalizeIdentity(identifier)
	return s.findUser(ctx, bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "username_normalized", Value: normalized}},
		bson.D{{Key: "email_normalized", Value: normalized}},
	}}})
}

// FindUserByID 只接受合法 ObjectID，防止认证上下文中的不透明字符串下沉为宽泛查询。
func (s *Store) FindUserByID(ctx context.Context, id string) (model.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, ErrUserNotFound
	}
	return s.findUser(ctx, bson.D{{Key: "_id", Value: objectID}})
}

// findUser 收敛单用户查询和 Driver 错误转换，保证所有公开查询返回相同的未找到语义。
func (s *Store) findUser(ctx context.Context, filter bson.D) (model.User, error) {
	var document userDocument
	if err := s.database.Collection("users").FindOne(ctx, filter).Decode(&document); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("查询用户: %w", err)
	}
	return userDocumentToModel(document), nil
}

// normalizeIdentity 只用于唯一查询键；展示用 username 仍保留用户输入的原始大小写。
func normalizeIdentity(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// userDocumentToModel 在 MongoDB 边界把 ObjectID 转为不透明字符串，并移除所有 BSON 类型。
func userDocumentToModel(document userDocument) model.User {
	return model.User{
		ID:           document.ID.Hex(),
		Username:     document.Username,
		Email:        document.Email,
		PasswordHash: document.PasswordHash,
		CreatedAt:    document.CreatedAt,
		UpdatedAt:    document.UpdatedAt,
	}
}
