# MongoDB 数据模型

> 本文档记录当前阶段需要的集合和关键字段。业务见 [`backend/BACKEND_REQUIREMENTS.md`](./backend/BACKEND_REQUIREMENTS.md)，接口见 [`API_CONTRACT.md`](./API_CONTRACT.md)。

## 1. 基本规则

- BSON 字段使用 `snake_case`，API JSON 使用 `camelCase`。
- 私有数据都保存 `user_id`。
- 时间保存为 UTC Date，页面日期按用户时区计算。
- 重量使用毫克整数，字段后缀 `_mg`。
- 热量使用千分之一千卡整数，字段后缀 `_milli_kcal`。
- 缺失营养字段使用 `null` 或缺失，不能写成真实零值。
- 饮食记录保存营养快照。

当前集合：

- `users`
- `user_profiles`
- `nutrition_targets`
- `dietary_preferences`
- `foods`
- `meal_records`
- `uploads`

当前不建立菜谱、AI、推荐或每日汇总集合。每日汇总和趋势直接从饮食记录计算。

## 2. `users`

```javascript
{
  _id: ObjectId,
  username: string,
  username_normalized: string,
  email: string,
  email_normalized: string,
  password_hash: string,
  created_at: Date,
  updated_at: Date
}
```

索引：

```text
unique(username_normalized)
unique(email_normalized)
```

密码只保存安全哈希，不保存明文和 JWT。

## 3. `user_profiles`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  birth_date: "YYYY-MM-DD" | null,
  sex: string | null,
  height_millimeters: int64 | null,
  weight_mg: int64 | null,
  time_zone: string,
  created_at: Date,
  updated_at: Date
}
```

索引：`unique(user_id)`。

`time_zone` 使用 IANA 名称，例如 `Asia/Shanghai`。

## 4. `nutrition_targets`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  targets: {
    calories_milli_kcal: int64 | null,
    protein_mg: int64 | null,
    fat_mg: int64 | null,
    carbohydrate_mg: int64 | null,
    fiber_mg: int64 | null
  },
  created_at: Date,
  updated_at: Date
}
```

索引：`unique(user_id)`。所有目标值非负。

## 5. `dietary_preferences`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  allergens: [string],
  avoided_foods: [string],
  notes: string,
  created_at: Date,
  updated_at: Date
}
```

索引：`unique(user_id)`。数组保存去除空白和重复后的文本。

## 6. `foods`

```javascript
{
  _id: ObjectId,
  name: string,
  name_normalized: string,
  aliases: [string],
  category: string,
  source_type: "seed" | "import" | "custom",
  source_key: string | null,
  owner_user_id: ObjectId | null,
  standard_serving_mg: int64 | null,
  nutrients_per_100g: {
    calories_milli_kcal: int64,
    protein_mg: int64,
    fat_mg: int64,
    carbohydrate_mg: int64,
    fiber_mg: int64 | null
  },
  created_at: Date,
  updated_at: Date,
  deleted_at: Date | null
}
```

来源规则：

| 类型 | 所有者 | 来源键 | 可见范围 |
| --- | --- | --- | --- |
| `seed` | `null` | 必填 | 所有用户 |
| `import` | `null` | 必填 | 所有用户 |
| `custom` | 当前用户 | `null` | 仅所有者 |

索引：

```text
unique(source_type, source_key)  # seed/import
(owner_user_id, deleted_at, name_normalized)
(source_type, deleted_at, name_normalized)
```

自定义 Food 删除时设置 `deleted_at`。历史饮食记录仍使用保存时快照。

## 7. `meal_records`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  occurred_at: Date,
  occurrence_time_zone: string,
  meal_type: "breakfast" | "lunch" | "dinner" | "snack",
  image_upload_id: ObjectId | null,
  items: [
    {
      food_id: ObjectId,
      food_name_snapshot: string,
      input_amount: {
        value_milli_units: int64,
        unit: "g" | "serving"
      },
      normalized_amount_mg: int64,
      nutrients_per_100g_snapshot: {
        calories_milli_kcal: int64,
        protein_mg: int64,
        fat_mg: int64,
        carbohydrate_mg: int64,
        fiber_mg: int64 | null
      },
      calculated_nutrients: {
        calories_milli_kcal: int64,
        protein_mg: int64,
        fat_mg: int64,
        carbohydrate_mg: int64,
        fiber_mg: int64 | null
      }
    }
  ],
  total_nutrients: {},
  data_complete: bool,
  created_at: Date,
  updated_at: Date,
  deleted_at: Date | null
}
```

计算公式：

```text
条目营养值 = 每 100g 营养值 × 实际重量 ÷ 100g
```

- 前端提交的营养合计不入库，后端重新计算。
- Food 修改或删除不改变已有快照。
- 删除饮食记录时设置 `deleted_at`，汇总和趋势排除它。

索引：

```text
(user_id, deleted_at, occurred_at desc)
(user_id, deleted_at, meal_type, occurred_at desc)
```

## 8. `uploads`

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  status: "pending" | "attached" | "deleted",
  storage_type: "local",
  storage_key: string,
  original_file_name: string,
  content_type: string,
  size_bytes: int64,
  meal_record_id: ObjectId | null,
  created_at: Date,
  updated_at: Date,
  deleted_at: Date | null
}
```

- `storage_key` 由服务端生成，不直接使用原文件名。
- 图片必须属于当前用户。
- 上传后为 `pending`，关联饮食记录后为 `attached`。
- 用户删除未关联图片或饮食记录移除图片时，删除文件并标记 `deleted`。

索引：

```text
(user_id, status, created_at desc)
unique(storage_key)
```

## 9. 每日汇总和趋势

当前不持久化汇总结果。

每日汇总查询当前用户指定自然日内未删除的 `meal_records`，累加 `total_nutrients`。最近 7 天重复相同规则，并由 Go 补齐没有记录的日期。

某项营养只要存在未知数据，该项汇总就标记为未知，不能把未知按零累加。

## 10. 简单关联

```text
users._id
  ├─ user_profiles.user_id
  ├─ nutrition_targets.user_id
  ├─ dietary_preferences.user_id
  ├─ foods.owner_user_id
  ├─ meal_records.user_id
  └─ uploads.user_id
```

MongoDB 不提供外键，Go Service 在写入前检查资源存在、所有者和状态。

## 11. Seed 和导入

Seed 与 JSON/CSV 导入至少包含：

- `source_key`
- `name`
- `category`
- 每 100g 热量、蛋白质、脂肪和碳水化合物
- 可选标准份量和纤维

使用 `source_type + source_key` 更新或插入，避免重复导入。

## 12. 变更规则

修改集合、字段、单位、索引或快照结构时同步更新本文档和 `API_CONTRACT.md`。小项目只为当前功能建模，不为未来 AI 或爬虫预建集合。
