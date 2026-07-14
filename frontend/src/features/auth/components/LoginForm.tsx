import { zodResolver } from '@hookform/resolvers/zod'
import { Controller, useForm } from 'react-hook-form'
import { useLocation, useNavigate } from 'react-router-dom'
import {
  Banner,
  Button,
  CheckboxInput,
  PasswordField,
  TextField,
  useAppToast,
} from '../../../components/ui'
import { ApiError, loginUser } from '../../../lib/api'
import { activateDemoSession } from '../../../lib/storage/session'
import { loginSchema, type LoginFormValues } from '../schemas/login-schema'

const defaultValues: LoginFormValues = {
  identifier: '',
  password: '',
  remember: false,
}

type RedirectLocation = {
  pathname?: string
  search?: string
  hash?: string
}

// 只接受应用内部地址，避免把路由状态当作任意跳转目标。
function getRedirectTarget(state: unknown): string {
  const from = (state as { from?: RedirectLocation } | null)?.from
  if (!from?.pathname?.startsWith('/app')) {
    return '/app/dashboard'
  }

  return `${from.pathname}${from.search ?? ''}${from.hash ?? ''}`
}

export function LoginForm() {
  const location = useLocation()
  const navigate = useNavigate()
  const showToast = useAppToast()
  const {
    control,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    defaultValues,
    mode: 'onBlur',
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = handleSubmit(async (values) => {
    try {
      const result = await loginUser({
        identifier: values.identifier.trim(),
        password: values.password,
      })

      activateDemoSession(result.sessionId, values.remember)
      showToast({
        body: `欢迎回来，${result.user.username}。`,
        autoHideDuration: 4000,
        uniqueID: 'login-success',
      })
      navigate(getRedirectTarget(location.state), { replace: true })
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        setError('root', { message: error.message })
        return
      }

      setError('root', {
        message: error instanceof Error ? error.message : '登录失败，请稍后重试。',
      })
    }
  })

  return (
    <form className="register-form" noValidate onSubmit={onSubmit}>
      {errors.root?.message && (
        <Banner
          description={errors.root.message}
          status="error"
          title="暂时无法登录"
        />
      )}

      <div className="register-form__fields">
        <Controller
          control={control}
          name="identifier"
          render={({ field }) => (
            <TextField
              {...field}
              error={errors.identifier?.message}
              hasAutoFocus
              htmlName={field.name}
              isRequired
              label="用户名或邮箱"
              placeholder="demo 或 demo@example.com"
              ref={field.ref}
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
              error={errors.password?.message}
              htmlName={field.name}
              isRequired
              label="密码"
              placeholder="输入登录密码"
              ref={field.ref}
              width="100%"
            />
          )}
        />
      </div>

      <Controller
        control={control}
        name="remember"
        render={({ field }) => (
          <CheckboxInput
            className="login-form__remember"
            htmlName={field.name}
            label="保持登录"
            onBlur={field.onBlur}
            onChange={(checked) => field.onChange(checked)}
            ref={field.ref}
            size="sm"
            value={field.value}
          />
        )}
      />

      <Button
        className="register-form__submit"
        isDisabled={isSubmitting}
        isLoading={isSubmitting}
        label={isSubmitting ? '正在登录' : '登录'}
        size="lg"
        type="submit"
        variant="primary"
      />
    </form>
  )
}
