export type NutritionSummary = {
  caloriesKcal: number
  proteinGrams: number
  fatGrams: number
  carbohydrateGrams: number
  fiberGrams?: number
}

export type UserProfile = {
  id: string
  username: string
  email: string
  birthDate?: string
  gender?: string
  heightCm?: number
  weightKg?: number
  healthRestrictions?: string[]
  dietaryAvoidances?: string[]
  notes?: string
}

// AuthUser 是认证接口允许返回的最小用户视图，不包含档案和任何敏感字段。
export type AuthUser = Pick<UserProfile, 'id' | 'username' | 'email'>

export type Food = {
  id: string
  name: string
  aliases: string[]
  category: string
  servingGrams: number
  nutritionPer100Grams: NutritionSummary
}

export type MealRecordItem = {
  foodId: string
  name: string
  weightGrams: number
  servings: number
  nutrition: NutritionSummary
}

export type MealRecord = {
  id: string
  eatenAt: string
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack'
  imageUrl?: string
  items: MealRecordItem[]
  nutrition: NutritionSummary
}

export type NutritionReport = {
  startDate: string
  endDate: string
  totals: NutritionSummary
  dailyAverage: NutritionSummary
  trend: Array<{ date: string; nutrition: NutritionSummary }>
}

export type Recommendation = {
  id: string
  foodId: string
  title: string
  reason: string
  suggestedServingGrams: number
  keyNutrients: string[]
  feedback?: 'liked' | 'disliked'
}

export type AppBootstrap = {
  user: AuthUser
}

export type RegisterRequest = {
  username: string
  email: string
  password: string
}

export type RegisterResponse = {
  accessToken: string
  tokenType: 'Bearer'
  expiresIn: number
  user: AuthUser
}

export type LoginRequest = {
  identifier: string
  password: string
}

export type LoginResponse = {
  accessToken: string
  tokenType: 'Bearer'
  expiresIn: number
  user: AuthUser
}
