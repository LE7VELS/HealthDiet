package store

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store 保存应用共享的 MongoDB Client 和数据库引用。
// 一个进程只创建一个 Store，以复用 Driver 内部连接池并集中控制关闭时机。
type Store struct {
	client   *mongo.Client
	database *mongo.Database
}

// New 连接 MongoDB，并在返回前 Ping 服务端；连接不可用时 Gin 不应启动。
func New(ctx context.Context, uri, databaseName string) (*Store, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("连接 MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		// Ping 失败时立即释放已创建的 Client，避免启动失败路径泄漏连接池资源。
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("Ping MongoDB: %w", err)
	}

	return &Store{
		client:   client,
		database: client.Database(databaseName),
	}, nil
}

// Close 在应用退出时关闭共享 Client，调用方应传入有限超时的 Context。
func (s *Store) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

// EnsureSchema 幂等补齐 DATA_MODEL.md 定义的集合和索引，不清空、迁移或写入虚构业务数据。
func (s *Store) EnsureSchema(ctx context.Context) error {
	// 显式创建集合可以尽早发现权限问题，并为后续索引初始化提供稳定目标。
	collections := []string{
		"users",
		"user_profiles",
		"nutrition_targets",
		"dietary_preferences",
		"foods",
		"meal_records",
		"uploads",
	}

	existing, err := s.database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("读取集合列表: %w", err)
	}

	// 集合列表转为集合结构，避免逐项创建前反复线性查找。
	existingSet := make(map[string]struct{}, len(existing))
	for _, name := range existing {
		existingSet[name] = struct{}{}
	}

	for _, name := range collections {
		if _, ok := existingSet[name]; ok {
			continue
		}

		// 只补齐缺失集合，不删除或清空已有业务数据。
		if err := s.database.CreateCollection(ctx, name); err != nil {
			return fmt.Errorf("创建集合 %s: %w", name, err)
		}
	}

	if err := s.ensureIndexes(ctx); err != nil {
		return err
	}

	return nil
}

// ensureIndexes 集中维护数据模型所需索引；命名索引允许应用重启时安全重复执行 CreateMany。
func (s *Store) ensureIndexes(ctx context.Context) error {
	indexes := map[string][]mongo.IndexModel{
		// 用户名和邮箱规范化字段必须唯一，这是并发注册冲突的最终一致性边界。
		"users": {
			uniqueIndex("username_normalized", bson.D{{Key: "username_normalized", Value: 1}}),
			uniqueIndex("email_normalized", bson.D{{Key: "email_normalized", Value: 1}}),
		},
		"user_profiles": {
			uniqueIndex("user_id", bson.D{{Key: "user_id", Value: 1}}),
		},
		"nutrition_targets": {
			uniqueIndex("user_id", bson.D{{Key: "user_id", Value: 1}}),
		},
		"dietary_preferences": {
			uniqueIndex("user_id", bson.D{{Key: "user_id", Value: 1}}),
		},
		"foods": {
			{
				Keys: bson.D{
					{Key: "source_type", Value: 1},
					{Key: "source_key", Value: 1},
				},
				Options: options.Index().
					SetName("source_type_source_key_unique").
					SetUnique(true).
					// 自定义 Food 的 source_key 为空，唯一约束只应用于 Seed 和导入数据。
					SetPartialFilterExpression(bson.D{
						{Key: "source_type", Value: bson.D{{Key: "$in", Value: bson.A{"seed", "import"}}}},
						{Key: "source_key", Value: bson.D{{Key: "$type", Value: "string"}}},
					}),
			},
			{
				Keys: bson.D{
					{Key: "owner_user_id", Value: 1},
					{Key: "deleted_at", Value: 1},
					{Key: "name_normalized", Value: 1},
				},
				Options: options.Index().SetName("owner_deleted_name"),
			},
			{
				Keys: bson.D{
					{Key: "source_type", Value: 1},
					{Key: "deleted_at", Value: 1},
					{Key: "name_normalized", Value: 1},
				},
				Options: options.Index().SetName("source_deleted_name"),
			},
		},
		// 饮食记录索引同时包含 user_id 与 deleted_at，支持用户隔离和软删除过滤。
		"meal_records": {
			{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
					{Key: "deleted_at", Value: 1},
					{Key: "occurred_at", Value: -1},
				},
				Options: options.Index().SetName("user_deleted_occurred"),
			},
			{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
					{Key: "deleted_at", Value: 1},
					{Key: "meal_type", Value: 1},
					{Key: "occurred_at", Value: -1},
				},
				Options: options.Index().SetName("user_deleted_meal_occurred"),
			},
		},
		"uploads": {
			{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
					{Key: "status", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("user_status_created"),
			},
			uniqueIndex("storage_key", bson.D{{Key: "storage_key", Value: 1}}),
		},
	}

	for collection, models := range indexes {
		// CreateMany 对相同定义的命名索引可重复执行，因此应用重启不会重复建索引。
		if _, err := s.database.Collection(collection).Indexes().CreateMany(ctx, models); err != nil {
			return fmt.Errorf("创建集合 %s 的索引: %w", collection, err)
		}
	}

	return nil
}

// uniqueIndex 统一命名简单唯一索引，降低重启时因自动生成名称不同而重复建索引的风险。
func uniqueIndex(name string, keys bson.D) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetName(name + "_unique").SetUnique(true),
	}
}
