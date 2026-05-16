import { MonthRangePicker } from './components/ui/month-picker'
import { format } from 'date-fns/format'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { AppSidebar } from '@/components/app-sidebar'
import { Switch } from './components/ui/switch'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription
} from '@/components/ui/field'
import { Button } from './components/ui/button'
import { SpinnerGapIcon } from '@phosphor-icons/react'
import { useForm, useStore } from '@tanstack/react-form'
import { Card } from './components/ui/card'
import { app, domain } from 'wails/go/models'
import { GetUser, RunWorkflow, StopWorkflow } from 'wails/go/main/App'
import { Tabs, TabsContent, TabsList, TabsTrigger } from './components/ui/tabs'
import { TerminalLogs } from './components/terminal'

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

      const date = new domain.DateRange()
      date.from = value.dateRange.from
      date.to = value.dateRange.to

      const work = new app.WorkFlowInput()
      work.user = await GetUser()
      work.date = date

      console.log(work)

      await RunWorkflow(work)
    }
  })

  const dates = useStore(form.store, state => state.values.dateRange)

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
          <Card className="py-0">
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
                    className="p-0 pb-4"
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
                            className="my-auto"
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
                      {isSubmitting ? (
                        <Button
                          type="button"
                          onClick={async () => await StopWorkflow()}
                        >
                          Cancel <SpinnerGapIcon className=" animate-spin" />
                        </Button>
                      ) : (
                        <Button type="submit">Submit</Button>
                      )}

                      <Button
                        disabled={!canSubmit}
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

          <Card className="min-h-screen flex-1 md:min-h-min">
            <Tabs defaultValue="terminal" className=" w-full">
              <TabsList>
                <TabsTrigger value="terminal">Terminal</TabsTrigger>
                <TabsTrigger value="analytics">Analytics</TabsTrigger>
              </TabsList>

              <TabsContent value="terminal">
                <TerminalLogs
                  initialEntries={[
                    {
                      id: 1,
                      time: new Date(Date.now() - 8000).toISOString(),
                      level: 'info' as const,
                      message: 'Server started',
                      fields: { port: 8080, env: 'development' }
                    },
                    {
                      id: 2,
                      time: new Date(Date.now() - 6000).toISOString(),
                      level: 'debug' as const,
                      message: 'Loading configuration from config.yaml'
                    },
                    {
                      id: 3,
                      time: new Date(Date.now() - 5000).toISOString(),
                      level: 'info' as const,
                      message: 'Connected to database',
                      fields: { host: 'localhost', db: 'myapp' }
                    },
                    {
                      id: 4,
                      time: new Date(Date.now() - 4000).toISOString(),
                      level: 'warning' as const,
                      message: 'Rate limit threshold at 80%',
                      fields: { limit: 1000, current: 802 }
                    },
                    {
                      id: 5,
                      time: new Date(Date.now() - 3000).toISOString(),
                      level: 'debug' as const,
                      message: 'Processing request GET /api/users'
                    },
                    {
                      id: 6,
                      time: new Date(Date.now() - 2000).toISOString(),
                      level: 'error' as const,
                      message: 'Failed to fetch remote config',
                      fields: { url: 'https://config.svc', status: 503 }
                    },
                    {
                      id: 7,
                      time: new Date(Date.now() - 1000).toISOString(),
                      level: 'info' as const,
                      message: 'Retry attempt 1/3'
                    }
                  ]}
                />
              </TabsContent>
            </Tabs>
          </Card>
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
