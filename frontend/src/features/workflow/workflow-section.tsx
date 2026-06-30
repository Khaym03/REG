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
import { LoadingButton } from '@/components/ui/loading-button'

export default function App() {
  const { workflowForm } = useAppForms()
  const dates = useStore(workflowForm.store, state => state.values.dateRange)
  const stopWorkflow = useWorkflowStore(state => state.stopWorkflow)
  const isDebouncing = useWorkflowStore(state => state.isDebouncing)

  return (
    <section className="@container size-full flex flex-col min-h-0 gap-4">
      <Card className="p-0 flex-none ring-0 ">
        <form
          onSubmit={e => {
            e.preventDefault()
            e.stopPropagation()
            workflowForm.handleSubmit()
          }}
          className="flex flex-col @min-3xl:flex-row justify-between gap-4"
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
                <div className="grid grid-cols-2 @min-3xl:flex @min-3xl:flex-col justify-end gap-1">
                  {isSubmitting ? (
                    <LoadingButton
                      type="button"
                      isLoading={isSubmitting}
                      loadingText="Cancelando..."
                      onClick={() => stopWorkflow()}
                    >
                      Cancel
                    </LoadingButton>
                  ) : (
                    <LoadingButton type="submit" isLoading={isSubmitting}>
                      Submit
                    </LoadingButton>
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

      <Card className="min-h-0 flex-1 flex flex-col p-0 ring-0  w-full">
        <StateFlow />
      </Card>
    </section>
  )
}
