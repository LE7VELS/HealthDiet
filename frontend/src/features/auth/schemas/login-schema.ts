import { z } from '../../../lib/validation'

export const loginSchema = z.object({
  identifier: z
    .string()
    .trim()
    .min(1, '请输入用户名或邮箱')
    .refine(
      (value) => !value.includes('@') || z.email().safeParse(value).success,
      '请输入有效的邮箱地址',
    ),
  password: z
    .string()
    .min(1, '请输入密码')
    .min(8, '密码至少需要 8 个字符')
    .max(64, '密码不能超过 64 个字符'),
  remember: z.boolean(),
})

export type LoginFormValues = z.infer<typeof loginSchema>
