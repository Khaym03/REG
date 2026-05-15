import { MonthRangePicker } from './components/ui/month-picker'
import { format } from 'date-fns/format'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { AppSidebar } from '@/components/app-sidebar'
import { Switch } from './components/ui/switch'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription,
} from '@/components/ui/field'
import { Button } from './components/ui/button'
import { useForm, useStore } from '@tanstack/react-form'
import { Card } from './components/ui/card'

interface WorkflowInput {
  dateRange: DateRange
  receptionPendingGuides: boolean
}

interface DateRange {
  from: Date
  to: Date
}

const defaultWorkflowInput: WorkflowInput = {
  dateRange: {
    from: new Date(),
    to: new Date()
  },
  receptionPendingGuides: false
}

function App() {
  const form = useForm({
    defaultValues: defaultWorkflowInput,
    onSubmit: async ({ value }) => {
      console.log(value)
    }
  })

  const dates = useStore(form.store, state => state.values.dateRange)

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
          <Card className='py-0'>
            <form
              onSubmit={e => {
                e.preventDefault()
                e.stopPropagation()
                form.handleSubmit()
              }}
              className="p-4 rounded-xl"
            >
              <form.Field
                name="dateRange"
                children={field => (
                  <MonthRangePicker
                    key={`${dates.from.toISOString()}-${dates.to.toISOString()}`}
                    onMonthRangeSelect={newDates => {
                      field.handleChange({
                        from: newDates.start,
                        to: newDates.end
                      })
                    }}
                    selectedMonthRange={{ start: dates.from, end: dates.to }}
                    maxDate={new Date()}
                    className='p-0 pb-4'
                  />
                )}
              />
              <div className="grid grid-cols-3 gap-2">
                <DisplaySelectedDate dates={dates} />
                <form.Field
                  name="receptionPendingGuides"
                  children={field => {
                    return (
                      <FieldLabel htmlFor={field.name}>
                        <Field orientation={'horizontal'}>
                          <FieldContent>
                            Include-pending-guides
                            <FieldDescription>
                              Pending guides will also be reception
                            </FieldDescription>
                          </FieldContent>
                          <Switch
                            id={field.name}
                            name={field.name}
                            checked={field.state.value}
                            onCheckedChange={field.handleChange}
                            className='my-auto'
                          />
                        </Field>
                      </FieldLabel>
                    )
                  }}
                />

                <form.Subscribe
                  selector={state => [state.canSubmit, state.isSubmitting]}
                  children={([canSubmit, isSubmitting]) => (
                    <div className="flex flex-col">
                      <Button type="submit" disabled={!canSubmit}>
                        {isSubmitting ? '...' : 'Submit'}
                      </Button>
                      <Button
                        variant={'secondary'}
                        type="reset"
                        onClick={e => {
                          // Avoid unexpected resets of form elements (especially <select> elements)
                          e.preventDefault()
                          form.reset()
                        }}
                      >
                        Reset
                      </Button>
                    </div>
                  )}
                />
              </div>
            </form>
          </Card>
          <Card className="min-h-screen flex-1 md:min-h-min" />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}

export default App

function DisplaySelectedDate({ dates }: { dates: DateRange }) {
  const startDateStr = format(dates.from, 'MMM yyyy')
  const endDateStr = format(dates.to, 'MMM yyyy')
  return (
    <FieldLabel>
      <Field>
        <FieldContent>
          Selected date:
         <FieldDescription>
           {startDateStr} - {endDateStr}
         </FieldDescription>
        </FieldContent>
      </Field>
    </FieldLabel>
  )
}
