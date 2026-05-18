import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription
} from '@/components/ui/field'
import { Switch } from './ui/switch'
import { useAppForms } from './use-app'

export default function SettingsSection() {
  const { browserForm } = useAppForms()
  return (
    <Card>
      <CardHeader>
        <CardTitle>Browser configuration</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={e => {
            e.preventDefault()
            e.stopPropagation()
            browserForm.handleSubmit()
          }}
        >
          <browserForm.Field
            name="headless"
            children={field => (
              <FieldLabel htmlFor={field.name}>
                <Field orientation={'horizontal'}>
                  <FieldContent>
                    Headless
                    <FieldDescription>
                      Execute the workflow in headless mode
                    </FieldDescription>
                  </FieldContent>
                  <Switch
                    id={field.name}
                    name={field.name}
                    checked={field.state.value}
                    onCheckedChange={field.handleChange}
                    className="my-auto"
                  />
                </Field>
              </FieldLabel>
            )}
          ></browserForm.Field>
          <browserForm.Field
            name="trace"
            children={field => (
              <FieldLabel htmlFor={field.name}>
                <Field orientation={'horizontal'}>
                  <FieldContent>
                    Trace
                    <FieldDescription>Enable browser trace</FieldDescription>
                  </FieldContent>
                  <Switch
                    id={field.name}
                    name={field.name}
                    checked={field.state.value}
                    onCheckedChange={field.handleChange}
                    className="my-auto"
                  />
                </Field>
              </FieldLabel>
            )}
          />
        </form>
      </CardContent>
    </Card>
  )
}
