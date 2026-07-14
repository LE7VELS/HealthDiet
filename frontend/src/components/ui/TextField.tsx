import {
  TextInput as AstryxTextInput,
  type TextInputProps,
} from '@astryxdesign/core/TextInput'
import { InputGroup } from '@astryxdesign/core/InputGroup'

export type TextFieldProps = Omit<TextInputProps, 'status'> & {
  error?: string
}

export function TextField({
  error,
  label,
  isRequired,
  isOptional,
  description,
  width,
  ...props
}: TextFieldProps) {
  // 统一使用 InputGroup，让普通输入框与密码框采用相同的独立错误提示布局。
  return (
    <InputGroup
      description={description}
      isOptional={isOptional}
      isRequired={isRequired}
      label={label}
      status={error ? { type: 'error', message: error } : undefined}
      style={{ width: typeof width === 'number' ? `${width}px` : width }}
    >
      <AstryxTextInput
        {...props}
        isLabelHidden
        isRequired={isRequired}
        label={label}
      />
    </InputGroup>
  )
}
