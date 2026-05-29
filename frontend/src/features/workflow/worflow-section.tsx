import { MonthRangePicker } from '@/components/ui/month-picker'

import { Switch } from '@/components/ui/switch'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription
} from '@/components/ui/field'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { useStore } from '@tanstack/react-form'
import { Card } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { TerminalLogs } from '@/features/workflow/components/terminal'
import DisplaySelectedDate from '@/features/workflow/components/display-selected-date'
import { useAppForms } from '../../hooks/use-app'
import StateFlow from './components/state-flow'
import { useWorkflowStore } from './store'
import { EventsOn } from 'wails/runtime/runtime'
import { useEffect } from 'react'

export default function App() {
  const { workflowForm } = useAppForms()
  const dates = useStore(workflowForm.store, state => state.values.dateRange)
  const stopWorkflow = useWorkflowStore(state => state.stopWorkflow)
  const isDebouncing = useWorkflowStore(state => state.isDebouncing)

  const addLogLine = useWorkflowStore(state => state.addLogLine)

  useEffect(() => {
    let unsubscribe: (() => void) | undefined

    const setupListener = async () => {
      try {
        unsubscribe = EventsOn('LOG', (line: string) => {
          addLogLine(line, 500)
        })
      } catch (err) {
        console.error('Failed to attach log listener', err)
      }
    }

    setupListener()
    return () => unsubscribe?.()
  }, [addLogLine])

  return (
    <>
      <Card className="p-0 flex-none ring-0">
        <form
          onSubmit={e => {
            e.preventDefault()
            e.stopPropagation()
            workflowForm.handleSubmit()
          }}
          className="flex justify-between gap-4"
        >
          <workflowForm.Field
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
            <workflowForm.Field
              name="receive_guides_in_transit"
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

            <workflowForm.Subscribe
              selector={state => [state.canSubmit, state.isSubmitting]}
              children={([canSubmit, isSubmitting]) => (
                <div className="flex flex-col justify-end gap-1">
                  {isSubmitting ? (
                    <Button
                      type="button"
                      disabled={isDebouncing}
                      onClick={() => {
                        stopWorkflow()
                      }}
                    >
                      <Spinner />
                      Cancel
                    </Button>
                  ) : (
                    <Button type="submit" disabled={isDebouncing}>
                      Submit
                    </Button>
                  )}

                  <Button
                    disabled={!canSubmit || isDebouncing}
                    variant={'secondary'}
                    type="reset"
                    onClick={e => {
                      // Avoid unexpected resets of form elements (especially <select> elements)
                      e.preventDefault()
                      workflowForm.reset()
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

      <Card className="min-h-0 flex flex-col p-0 ring-1">
        <Tabs
          defaultValue="terminal"
          className="w-full gap-0"
          style={{
            height: `300px`
          }}
        >
          <TabsList className="">
            <TabsTrigger value="terminal">Terminal</TabsTrigger>
            <TabsTrigger value="workflow">Workflow</TabsTrigger>
          </TabsList>

          <TabsContent
            value="terminal"
            className="flex-1 min-h-0 overflow-y-auto border overflow-x-hidden"
          >
            <TerminalLogs />
          </TabsContent>
          <TabsContent value="workflow">
            <StateFlow />
          </TabsContent>
        </Tabs>
      </Card>
    </>
  )
}
