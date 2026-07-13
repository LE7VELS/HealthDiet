import { z } from '../../../lib/validation'

const passwordSchema = z
  .string()
  .min(8, '密码至少需要 8 个字符')
  .max(64, '密码不能超过 64 个字符')
  .regex(/[A-Za-z]/, '密码至少需要包含一个英文字母')
  .regex(/\d/, '密码至少需要包含一个数字')

export const registerSchema = z
  .object({
    username: z
      .string()
      .trim()
      .min(2, '用户名至少需要 2 个字符')
      .max(24, '用户名不能超过 24 个字符')
      .regex(/^[\p{L}\p{N}_-]+$/u, '用户名只能包含文字、数字、下划线或连字符'),
    email: z.string().trim().min(1, '请输入邮箱').email('请输入有效的邮箱地址'),
    password: passwordSchema,
    confirmPassword: z.string().min(1, '请再次输入密码'),
  })
  .refine((values) => values.password === values.confirmPassword, {
    message: '两次输入的密码不一致',
    path: ['confirmPassword'],
  })

export type RegisterFormValues = z.infer<typeof registerSchema>
