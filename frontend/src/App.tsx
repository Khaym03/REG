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
import { SpinnerGapIcon, XIcon } from '@phosphor-icons/react'
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
    <div className="overflow-hidden">
      <SidebarProvider>
        <SidebarInset>
          <div
            style={{ '--wails-draggable': 'drag' }}
            className="h-7 w-full bg-background flex justify-end fixed top-0 right-0 z-0"
          ></div>
          <div className="relative flex flex-1 flex-col p-4 py-0 overflow-y-hidden mt-7 h-[572px] max-h-[572px] gap-4">
            <Card className="p-0 flex-none ring-0">
              <form
                onSubmit={e => {
                  e.preventDefault()
                  e.stopPropagation()
                  form.handleSubmit()
                }}
                className="flex justify-between gap-4"
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
                      className="p-0 flex-1"
                    />
                  )}
                />
                <div className="grid grid-cols-1 gap-2">
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
                      <div className="flex flex-col justify-end gap-1">
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

            <Card className="min-h-0 flex flex-col p-0 ring-0">
              <Tabs
                defaultValue="terminal"
                className="w-full gap-0"
                style={{
                  height: `300px`
                }}
              >
                <TabsList className="">
                  <TabsTrigger value="terminal">Terminal</TabsTrigger>
                  <TabsTrigger value="analytics">Analytics</TabsTrigger>
                </TabsList>

                <TabsContent
                  value="terminal"
                  className="flex-1 min-h-0 overflow-y-auto border"
                >
                  <TerminalLogs />
                </TabsContent>
              </Tabs>
            </Card>
          </div>
        </SidebarInset>
        <AppSidebar side="right" />
      </SidebarProvider>
    </div>
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
