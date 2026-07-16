import { zodResolver } from '@hookform/resolvers/zod'
import { Controller, useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { Banner, Button, PasswordField, Text, TextField, useAppToast } from '../../../components/ui'
import { ApiError, registerUser } from '../../../lib/api'
import { activateSession } from '../../../lib/storage/session'
import { registerSchema, type RegisterFormValues } from '../schemas/register-schema'

const defaultValues: RegisterFormValues = {
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
}

export function RegisterForm() {
  const navigate = useNavigate()
  const showToast = useAppToast()
  const {
    control,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormValues>({
    defaultValues,
    mode: 'onBlur',
    resolver: zodResolver(registerSchema),
  })

  const onSubmit = handleSubmit(async (values) => {
    try {
      const result = await registerUser({
        username: values.username.trim(),
        email: values.email.trim().toLowerCase(),
        password: values.password,
      })

      activateSession(result)
      showToast({
        body: '注册成功，接下来可以完善个人档案。',
        autoHideDuration: 2000,
        uniqueID: 'registration-success',
      })
      navigate('/app/profile', { replace: true })
    } catch (error) {
      if (error instanceof ApiError && error.code === 'VALIDATION_ERROR' && error.fields.length > 0) {
        // 服务端校验是最终边界，将合同中的字段错误落回对应输入框并保留用户填写内容。
        for (const fieldError of error.fields) {
          if (fieldError.field === 'username' || fieldError.field === 'email' || fieldError.field === 'password') {
            setError(fieldError.field, { message: fieldError.message })
          }
        }
        return
      }
      if (error instanceof ApiError && error.code === 'USERNAME_CONFLICT') {
        setError('username', { message: error.message }, { shouldFocus: true })
        return
      }
      if (error instanceof ApiError && error.code === 'EMAIL_CONFLICT') {
        setError('email', { message: error.message }, { shouldFocus: true })
        return
      }

      setError('root', {
        message: error instanceof Error ? error.message : '注册失败，请稍后重试。',
      })
    }
  })

  return (
    <form className="register-form" noValidate onSubmit={onSubmit}>
      {errors.root?.message && (
        <Banner
          description={errors.root.message}
          status="error"
          title="暂时无法注册"
        />
      )}

      <div className="register-form__fields">
        <Controller
          control={control}
          name="username"
          render={({ field }) => (
            <TextField
              {...field}
              error={errors.username?.message}
              hasAutoFocus
              htmlName={field.name}
              isRequired
              label="用户名"
              placeholder="例如：健康生活家"
              ref={field.ref}
              width="100%"
            />
          )}
        />

        <Controller
          control={control}
          name="email"
          render={({ field }) => (
            <TextField
              {...field}
              error={errors.email?.message}
              htmlName={field.name}
              isRequired
              label="邮箱"
              placeholder="name@example.com"
              ref={field.ref}
              type="email"
              width="100%"
            />
          )}
        />

        <Controller
          control={control}
          name="password"
          render={({ field }) => (
            <PasswordField
              {...field}
              description="至少 8 个字符，并包含英文字母和数字。"
              error={errors.password?.message}
              htmlName={field.name}
              isRequired
              label="密码"
              placeholder="设置登录密码"
              ref={field.ref}
              width="100%"
            />
          )}
        />

        <Controller
          control={control}
          name="confirmPassword"
          render={({ field }) => (
            <PasswordField
              {...field}
              error={errors.confirmPassword?.message}
              htmlName={field.name}
              isRequired
              label="确认密码"
              placeholder="再次输入密码"
              ref={field.ref}
              width="100%"
            />
          )}
        />
      </div>

      <Text as="p" className="register-form__notice" color="secondary" type="supporting">
        注册即表示你理解本产品仅提供日常膳食记录与营养参考，不替代医疗诊断或专业建议。
      </Text>

      <Button
        className="register-form__submit"
        isDisabled={isSubmitting}
        isLoading={isSubmitting}
        label={isSubmitting ? '正在创建账号' : '创建账号'}
        size="lg"
        type="submit"
        variant="primary"
      />
    </form>
  )
}
