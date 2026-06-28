import { MonthRangePicker } from '@/components/ui/month-picker'

import { Switch } from '@/components/ui/switch'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription
} from '@/components/ui/field'
import { Button } from '@/components/ui/button'
import { useStore } from '@tanstack/react-form'
import { Card } from '@/components/ui/card'
import DisplaySelectedDate from '@/features/workflow/components/display-selected-date'
import { useAppForms } from '../../hooks/use-app'
import StateFlow from './components/state-flow'
import { useWorkflowStore } from './store'

export default function App() {
  const { workflowForm } = useAppForms()
  const dates = useStore(workflowForm.store, state => state.values.dateRange)
  const stopWorkflow = useWorkflowStore(state => state.stopWorkflow)
  const isDebouncing = useWorkflowStore(state => state.isDebouncing)

  return (
    <section className='w-3xl'>
      <Card className="p-0 flex-none ring-0 ">
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
                      className="transition-all"
                    >
                      <div className="shiny inline-block bg-[linear-gradient(120deg,rgba(255,255,255,0)_40%,rgba(255,255,255,0.8)_50%,rgba(255,255,255,0)_60%)] dark:bg-[linear-gradient(120deg,rgba(0,0,0,0)_40%,rgba(0,0,0,0.8)_50%,rgba(0,0,0,0)_60%)] bg-size-[200%_100%] bg-clip-text text-white/70 dark:text-foreground/60">
                        Cancel
                      </div>
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

      <Card className="min-h-0 flex flex-col p-0 ring-1 h-[300px] w-full">
        <StateFlow />
      </Card>
    </section>
  )
}
