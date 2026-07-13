import {
  TextInput as AstryxTextInput,
  type TextInputProps,
} from '@astryxdesign/core/TextInput'

export type TextFieldProps = Omit<TextInputProps, 'status'> & {
  error?: string
}

export function TextField({ error, ...props }: TextFieldProps) {
  return (
    <AstryxTextInput
      {...props}
      status={error ? { type: 'error', message: error } : undefined}
    />
  )
}
