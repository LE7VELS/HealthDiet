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
- `internal/model` 使用不依赖 Gin 和 MongoDB Driver 的 Go 业务结构；Store 使用带 BSON 标签的内部 Document，并负责两者转换。

当前集合：

- `users`
- `user_profiles`
- `nutrition_targets`
- `dietary_preferences`
- `foods`
- `meal_records`
- `uploads`

当前不建立独立菜谱、AI、推荐或每日汇总集合。食品和菜谱统一保存在 `foods`；每日汇总和趋势直接从饮食记录计算。

Go 后端启动时幂等创建上述集合和本文件定义的索引。初始化只补齐结构，不清空集合，也不自动插入虚构 Food 数据。

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
  kind: "food" | "recipe",
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
  recipe: {
    ingredients: [
      {
        food_id: ObjectId,
        food_name_snapshot: string,
        amount_mg: int64,
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
    steps: [string],
    servings: int32,
    finished_weight_mg: int64 | null,
    calculation_weight_mg: int64,
    calculation_basis: "finished_weight" | "ingredient_weight_estimate",
    total_nutrients: {
      calories_milli_kcal: int64,
      protein_mg: int64,
      fat_mg: int64,
      carbohydrate_mg: int64,
      fiber_mg: int64 | null
    },
    per_serving_nutrients: {
      calories_milli_kcal: int64,
      protein_mg: int64,
      fat_mg: int64,
      carbohydrate_mg: int64,
      fiber_mg: int64 | null
    }
  } | null,
  created_at: Date,
  updated_at: Date,
  deleted_at: Date | null
}
```

- `kind=food` 时 `recipe` 为 `null`。
- `kind=recipe` 表示菜品/菜谱，当前数据模型不区分这两个中文概念；`ingredients`、`steps` 至少各有一项，`servings` 大于 0。
- 原料必须引用创建者可访问的 `kind=food`，当前不允许引用菜谱；原料用量使用毫克整数且大于 0。
- 创建或修改菜谱时保存原料名称和营养快照，基础食品后续修改不会悄悄改变已有菜谱；用户主动编辑菜谱时重新读取原料并计算。
- 菜谱总营养等于各原料 `calculated_nutrients` 之和，每份营养等于总营养除以 `servings`。
- 提供 `finished_weight_mg` 时，`calculation_weight_mg` 使用成品重量；否则使用原料重量之和，并将 `calculation_basis` 标记为估算。
- 顶层 `nutrients_per_100g` 由总营养和 `calculation_weight_mg` 计算；顶层 `standard_serving_mg` 由 `calculation_weight_mg / servings` 计算。
- 步骤文本当前不参与营养计算，不根据做法猜测营养损耗、吸油量或水分变化。
- 任一原料缺少热量或宏量营养时不保存菜谱；纤维缺失时，菜谱对应的纤维快照和计算结果保持 `null`。

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
      food_kind_snapshot: "food" | "recipe",
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

foods._id
  ├─ foods.recipe.ingredients.food_id
  └─ meal_records.items.food_id
```

MongoDB 不提供外键。复杂写入由 Go Service 检查资源存在、所有者和状态；简单私有查询由 Store 强制要求 `user_id`。

## 11. Seed 和导入

Seed 与 JSON/CSV 导入至少包含：

- `source_key`
- `name`
- `kind`
- `category`
- `kind=food` 时提供每 100g 热量、蛋白质、脂肪和碳水化合物，以及可选标准份量和纤维
- `kind=recipe` 时提供基础食品原料引用、克数、份数、可选成品重量和制作步骤；使用 JSON 导入嵌套结构，营养仍由后端计算

使用 `source_type + source_key` 更新或插入，避免重复导入。

## 12. 变更规则

修改集合、字段、单位、索引或快照结构时同步更新本文档和 `API_CONTRACT.md`。小项目只为当前功能建模，不为未来 AI 或爬虫预建集合。

## 13. 后续食品知识库和 AI 数据边界

本节只记录未来扩展方向，当前集合清单和实现范围保持不变。

- `users`、`user_profiles`、`nutrition_targets`、`dietary_preferences`、`meal_records` 和 `uploads` 属于核心业务数据，只能由 Go API 按用户权限访问。
- 公共食品资料、来源文档和后续检索数据应与核心业务数据隔离，可使用独立数据库、独立集合边界或只读视图；具体方案在实施知识库时确定。
- AI 如需直接检索公共食品知识，只能使用独立的最小权限只读凭证，不能读取用户私有集合或写入数据。
- 爬虫原始数据、清洗结果和正式 Food 应区分处理。未经清洗、去重和校验的数据不能直接进入正式 `foods`，也不能作为营养计算依据。
- Food 后续可以扩展来源地址、外部数据 ID、抓取或导入批次、更新时间、数据版本、验证状态和质量等级等来源元数据；字段在真实需求确定后再建模。
- 食品知识文本用于检索、问答和解释，结构化 Food 用于确定性营养计算，两者不能互相替代。
- 当前不创建知识库、爬虫、Embedding、向量、Agent 会话或推荐结果集合。
