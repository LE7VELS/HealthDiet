# API 合同

> 本文档只记录前后端实现需要统一的接口。业务范围见 [`PROJECT_CONTEXT.md`](./PROJECT_CONTEXT.md)。

## 1. 通用约定

- 基础路径：`/api/v1`。
- JSON 字段使用 `camelCase`。
- 普通请求使用 `application/json`，图片使用 `multipart/form-data`。
- 资源 ID 是不透明字符串，前端不依赖 MongoDB ObjectID 格式。
- 时间使用 RFC 3339，日期使用 `YYYY-MM-DD`，时区使用 IANA 名称。
- 热量单位为 `kcal`，蛋白质、脂肪、碳水化合物和纤维单位为 `g`。
- `null` 表示数据缺失，`0` 表示真实为零。

成功响应：

```json
{"data": {}}
```

列表响应：

```json
{
  "data": [],
  "pagination": {"page": 1, "pageSize": 20, "total": 0}
}
```

错误响应：

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "请求参数不正确",
    "fields": [
      {"field": "email", "message": "邮箱格式不正确"}
    ]
  }
}
```

常用状态码：`200`、`201`、`204`、`400`、`401`、`404`、`409`、`413`、`415`、`500`。

客户端以 `error.code` 进行稳定的分支处理，不解析 `message` 判断错误类型。`message` 是可直接展示的安全消息；`fields` 只在错误能够定位到具体请求字段时返回，字段名使用请求 JSON 的 `camelCase` 名称。服务端内部错误链、数据库错误和实现细节不属于 API 响应。

## 2. 认证

受保护接口使用：

```http
Authorization: Bearer <accessToken>
```

当前只定义访问 Token，不做刷新 Token。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/auth/register` | 注册并返回 Token |
| POST | `/auth/login` | 登录并返回 Token |
| GET | `/auth/me` | 当前账号 |

注册请求：

```json
{"username":"demo","email":"demo@example.com","password":"password"}
```

登录请求：

```json
{"identifier":"demo@example.com","password":"password"}
```

登录或注册响应：

```json
{
  "data": {
    "accessToken": "jwt",
    "tokenType": "Bearer",
    "expiresIn": 3600,
    "user": {"id":"id","username":"demo","email":"demo@example.com"}
  }
}
```

## 3. 用户档案

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/profile` | 获取档案 |
| PATCH | `/profile` | 更新档案 |

档案对象：

```json
{
  "username": "demo",
  "email": "demo@example.com",
  "birthDate": null,
  "sex": null,
  "heightCm": null,
  "weightKg": null,
  "timeZone": "Asia/Shanghai"
}
```

`username` 和 `email` 当前只读。可选字段传 `null` 表示清空。

## 4. 营养目标和饮食偏好

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/nutrition-target` | 获取目标 |
| PUT | `/nutrition-target` | 保存目标 |
| GET | `/dietary-preferences` | 获取偏好 |
| PUT | `/dietary-preferences` | 保存偏好 |

营养目标：

```json
{
  "caloriesKcal": 2000,
  "proteinG": 80,
  "fatG": 60,
  "carbohydrateG": 250,
  "fiberG": null
}
```

饮食偏好：

```json
{
  "allergens": ["花生"],
  "avoidedFoods": ["香菜"],
  "notes": ""
}
```

## 5. Food

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/foods` | 查询食品或菜谱 |
| GET | `/foods/{foodId}` | Food 详情 |
| POST | `/foods/custom` | 创建自定义 Food |
| PATCH | `/foods/{foodId}` | 修改自己的自定义 Food |
| DELETE | `/foods/{foodId}` | 删除自己的自定义 Food |

查询参数：`query`、`kind=food|recipe`、`category`、`scope=public|custom`、`page`、`pageSize`。默认每页 20，最大 100。

Food 对象：

```json
{
  "id": "food-id",
  "name": "燕麦片",
  "aliases": ["燕麦"],
  "kind": "food",
  "category": "staple",
  "sourceType": "seed",
  "isOwnedByCurrentUser": false,
  "standardServingGrams": 40,
  "nutrientsPer100g": {
    "caloriesKcal": 377,
    "proteinG": 15,
    "fatG": 6.7,
    "carbohydrateG": 66.3,
    "fiberG": 10.1
  },
  "recipe": null
}
```

`sourceType`：`seed`、`import`、`custom`。

`kind`：`food` 表示食品，`recipe` 表示菜品/菜谱。当前系统不区分“菜品”和“菜谱”；两种类型都是 Food，可以加入饮食记录。

`category`：`staple`、`meat_egg_dairy`、`vegetable`、`fruit`、`legume`、`dish`、`beverage`、`other`。

菜谱类型的 Food 响应示例：

```json
{
  "id": "recipe-id",
  "name": "番茄炒鸡蛋",
  "aliases": [],
  "kind": "recipe",
  "category": "dish",
  "sourceType": "custom",
  "isOwnedByCurrentUser": true,
  "standardServingGrams": 165,
  "nutrientsPer100g": {
    "caloriesKcal": 112.4,
    "proteinG": 6.1,
    "fatG": 7.2,
    "carbohydrateG": 4.9,
    "fiberG": 0.8
  },
  "recipe": {
    "ingredients": [
      {"foodId": "egg-id", "foodName": "鸡蛋", "amountGrams": 120},
      {"foodId": "tomato-id", "foodName": "番茄", "amountGrams": 200},
      {"foodId": "oil-id", "foodName": "食用油", "amountGrams": 10}
    ],
    "steps": ["番茄切块，鸡蛋打散", "炒熟鸡蛋后加入番茄"],
    "servings": 2,
    "finishedWeightGrams": null,
    "calculationBasis": "ingredient_weight_estimate",
    "totalNutrients": {
      "caloriesKcal": 370.9,
      "proteinG": 20.1,
      "fatG": 23.8,
      "carbohydrateG": 16.2,
      "fiberG": 2.6
    },
    "perServingNutrients": {
      "caloriesKcal": 185.5,
      "proteinG": 10.1,
      "fatG": 11.9,
      "carbohydrateG": 8.1,
      "fiberG": 1.3
    }
  }
}
```

`kind=food` 时 `recipe` 必须为 `null`。`kind=recipe` 时：

- `ingredients` 和 `steps` 至少各有一项，`servings` 必须大于 0。
- 每个原料通过 `foodId` 引用当前用户可访问的 `kind=food`，`amountGrams` 必须大于 0；当前不允许引用另一个菜谱。
- `finishedWeightGrams` 可选；填写时必须大于 0。
- `calculationBasis=finished_weight` 表示每 100g 营养按成品重量计算；`ingredient_weight_estimate` 表示未填写成品重量，暂按原料总重量估算。
- `foodName`、`totalNutrients`、`perServingNutrients`、顶层 `standardServingGrams` 和 `nutrientsPer100g` 均由后端计算，客户端不得提交或覆盖。
- 制作步骤当前只用于记录和展示，不参与营养公式。
- 任一原料缺少热量、蛋白质、脂肪或碳水化合物时返回 `FOOD_DATA_INCOMPLETE`；纤维缺失时相关菜谱纤维值返回 `null`。

创建普通食品时提交 `name`、`aliases`、`kind=food`、`category`、`standardServingGrams` 和 `nutrientsPer100g`。热量、蛋白质、脂肪和碳水化合物必填且非负；纤维和标准份量可以为 `null`。

创建菜谱时提交：

```json
{
  "name": "番茄炒鸡蛋",
  "aliases": [],
  "kind": "recipe",
  "category": "dish",
  "recipe": {
    "ingredients": [
      {"foodId": "egg-id", "amountGrams": 120},
      {"foodId": "tomato-id", "amountGrams": 200},
      {"foodId": "oil-id", "amountGrams": 10}
    ],
    "steps": ["番茄切块，鸡蛋打散", "炒熟鸡蛋后加入番茄"],
    "servings": 2,
    "finishedWeightGrams": null
  }
}
```

服务端生成 `id`，并固定 `sourceType=custom` 和当前用户所有权。修改普通食品时接受其可编辑字段的子集；修改菜谱组成、份数、成品重量或步骤时提交完整 `recipe`，由后端重新校验并计算。当前不支持在 `food` 和 `recipe` 之间直接修改 `kind`，需要新建正确类型的 Food。

## 6. 图片

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/uploads/images` | 上传图片 |
| GET | `/uploads/images/{uploadId}/content` | 读取自己的图片 |
| DELETE | `/uploads/images/{uploadId}` | 删除未关联图片 |

上传字段名为 `file`，只允许 JPEG、PNG、WebP，最大 10 MB。

上传响应：

```json
{
  "data": {
    "id": "upload-id",
    "contentType": "image/jpeg",
    "sizeBytes": 245000,
    "contentUrl": "/api/v1/uploads/images/upload-id/content",
    "status": "pending"
  }
}
```

创建或更新饮食记录时提交 `imageId` 完成关联。记录保存失败时，前端可以重试或删除这张未关联图片。

## 7. 饮食记录

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/meals` | 创建记录 |
| GET | `/meals` | 查询记录 |
| GET | `/meals/{mealId}` | 记录详情 |
| PATCH | `/meals/{mealId}` | 更新记录 |
| DELETE | `/meals/{mealId}` | 删除记录 |

创建请求：

```json
{
  "occurredAt": "2026-07-14T12:30:00+08:00",
  "mealType": "lunch",
  "imageId": null,
  "items": [
    {
      "foodId": "food-id",
      "amount": {"value": 80, "unit": "g"}
    }
  ]
}
```

`mealType`：`breakfast`、`lunch`、`dinner`、`snack`。

数量单位：

- `g`：克。
- `serving`：份，要求 Food 有 `standardServingGrams`。

记录对象：

```json
{
  "id": "meal-id",
  "occurredAt": "2026-07-14T04:30:00Z",
  "localDate": "2026-07-14",
  "mealType": "lunch",
  "image": null,
  "items": [
    {
      "foodId": "food-id",
      "foodName": "燕麦片",
      "foodKind": "food",
      "amount": {"value": 80, "unit": "g"},
      "normalizedGrams": 80,
      "nutrients": {
        "caloriesKcal": 301.6,
        "proteinG": 12,
        "fatG": 5.36,
        "carbohydrateG": 53.04,
        "fiberG": 8.08
      }
    }
  ],
  "totalNutrients": {},
  "dataComplete": true
}
```

列表查询参数：`dateFrom`、`dateTo`、`mealType`、`query`、`page`、`pageSize`。默认按用餐时间倒序。

更新时提交 `imageId: null` 表示移除图片；提交 `items` 时提交完整新数组，后端重新计算营养快照。

## 8. 每日汇总和 7 天趋势

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/nutrition/daily?date=2026-07-14` | 每日汇总 |
| GET | `/nutrition/trends/7-days?endDate=2026-07-14` | 最近 7 天 |

每日汇总：

```json
{
  "data": {
    "date": "2026-07-14",
    "recordCount": 3,
    "actual": {},
    "target": {},
    "dataComplete": true,
    "hints": [
      {
        "code": "CALORIES_BELOW_TARGET",
        "metric": "caloriesKcal",
        "status": "low",
        "message": "今日记录热量低于目标。"
      }
    ]
  }
}
```

趋势返回连续 7 个日期；没有记录的日期也返回一个零记录数据点。

规则提示只来自后端确定性规则，不包含 AI 内容。

## 9. 主要错误码

| 错误码 | 含义 |
| --- | --- |
| `VALIDATION_ERROR` | 输入错误 |
| `INVALID_CREDENTIALS` | 登录失败 |
| `UNAUTHENTICATED` | 未登录或 Token 无效 |
| `RESOURCE_NOT_FOUND` | 当前用户范围内资源不存在 |
| `USERNAME_CONFLICT` | 用户名重复 |
| `EMAIL_CONFLICT` | 邮箱重复 |
| `FOOD_DATA_INCOMPLETE` | Food 数据不足以计算 |
| `UNSUPPORTED_IMAGE_TYPE` | 图片类型不支持 |
| `IMAGE_TOO_LARGE` | 图片过大 |
| `INTERNAL_ERROR` | 内部错误 |

## 10. 变更规则

修改接口路径、方法、字段、单位、枚举或错误码时，同步更新前端 DTO、后端实现和本文档。

根路径 `GET /` 只用于开发时手工确认 Gin 已启动，不属于业务 API。

## 11. 后续 Agent 接口扩展原则

本节记录未来方向，不属于当前 API 合同，也不要求当前实现。

- 未来 Agent 仍通过版本化、受认证的 Go API 读取用户私有数据和执行写操作。
- 可以增加聚合的 Agent 上下文接口，一次返回档案、营养目标、饮食偏好、近期记录、每日汇总和趋势，减少多次工具调用。
- Agent 可以直接只读检索与用户业务库隔离的公共食品知识库；知识库凭证不得用于访问用户集合。
- 写接口不得接收任意 MongoDB 查询或更新表达式。Agent 只提交结构化业务请求，Go Service 完成权限、数据和营养规则校验。
- 涉及新增或修改用户数据时，需要用户确认；具体路径、请求字段和鉴权方式在进入 AI 阶段后再写入正式合同。
