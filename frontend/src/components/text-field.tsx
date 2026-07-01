import { Input } from '@/components/ui/input'
import { Field, FieldDescription, FieldError, FieldLabel } from './ui/field'
import { useFieldContext } from '@/hooks/form-context'
import { useStore } from '@tanstack/react-form'

interface FormTextFieldProps {
  label: string
  placeholder?: string
  type?: 'text' | 'password' | 'email'
  description?: string
  required?: boolean
}

export default function TextField({
  label,
  placeholder,
  type = 'text',
  description,
  required = false
}: FormTextFieldProps) {
  const field = useFieldContext<string>()
  const errors = useStore(field.store, state => state.meta.errors)
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid

  return (
    <Field>
      <FieldLabel htmlFor={field.name}>{label}</FieldLabel>

      <Input
        id={field.name}
        name={field.name}
        type={type}
        placeholder={placeholder}
        value={field.state.value}
        onBlur={field.handleBlur}
        onChange={e => field.handleChange(e.target.value)}
        aria-invalid={isInvalid}
        required={required}
      />

      {description && <FieldDescription>{description}</FieldDescription>}
      {isInvalid && <FieldError errors={errors} />}
    </Field>
  )
}
