import { useState } from 'react'
import { InputGroup, InputGroupText } from '@astryxdesign/core/InputGroup'
import { TextInput, type TextInputProps } from '@astryxdesign/core/TextInput'

type PasswordFieldProps = Omit<TextInputProps, 'type' | 'status'> & {
  error?: string
}

export function PasswordField({
  error,
  label,
  isRequired,
  description,
  ...props
}: PasswordFieldProps) {
  const [isVisible, setIsVisible] = useState(false)

  return (
    <InputGroup
      description={description}
      isRequired={isRequired}
      label={label}
      status={error ? { type: 'error', message: error } : undefined}
    >
      <TextInput
        {...props}
        isLabelHidden
        isRequired={isRequired}
        label={label}
        type={isVisible ? 'text' : 'password'}
      />
      <InputGroupText>
        <button
          aria-label={`${isVisible ? '隐藏' : '显示'}${label}`}
          aria-pressed={isVisible}
          className="ui-password-toggle"
          onClick={() => setIsVisible((value) => !value)}
          type="button"
        >
          {isVisible ? '隐藏' : '显示'}
        </button>
      </InputGroupText>
    </InputGroup>
  )
}
