import { MonthRangePicker } from '@/components/ui/month-picker'

import { Switch } from '@/components/ui/switch'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription
} from '@/components/ui/field'
import { Button } from '@/components/ui/button'
import { SpinnerGapIcon } from '@phosphor-icons/react'
import { useStore } from '@tanstack/react-form'
import { Card } from '@/components/ui/card'
import { StopWorkflow } from 'wails/go/main/App'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { TerminalLogs } from '@/components/terminal'
import DisplaySelectedDate from '@/components/display-selected-date'
import { useAppForms } from './use-app'
import StateFlow from './state-flow'

export default function App() {
  const { workflowForm } = useAppForms()
  const dates = useStore(workflowForm.store, state => state.values.dateRange)

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
            <TabsTrigger value="workflow">Workflow</TabsTrigger>
          </TabsList>

          <TabsContent
            value="terminal"
            className="flex-1 min-h-0 overflow-y-auto border"
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
