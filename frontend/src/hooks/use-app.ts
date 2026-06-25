import { createContext, useContext } from 'react'
import { useForm } from '@tanstack/react-form'
// import { workflow, config, domain } from "bindings/github.com/Khaym03/REG";
import * as workflow from 'bindings/github.com/Khaym03/REG/internal/workflow'
import * as domain from 'bindings/github.com/Khaym03/REG/internal/domain'
import * as config from 'bindings/github.com/Khaym03/REG/internal/config'
import { App } from 'bindings/github.com/Khaym03/REG'
import type { WorkflowInput } from '@/types/types'
import { useWorkflowStore } from '@/features/workflow/store'

export type BrowserConfigForm = ReturnType<typeof useBrowserConfigFormInstance>
const defaultBrowserConfig = new config.BrowserConfig({
  headless: true,
  trace: true
})
export function useBrowserConfigFormInstance() {
  return useForm({
    defaultValues: defaultBrowserConfig,
    onSubmit: async ({ value }) => {
      console.log(value)
    }
  })
}

export type WorkflowForm = ReturnType<typeof useWorkflowFormInstance>
const defaultWorkflowInput: WorkflowInput = {
  dateRange: {
    from: new Date(),
    to: new Date()
  },
  receive_guides_in_transit: false
}
export function useWorkflowFormInstance(browserForm: BrowserConfigForm) {
  const runWorkflow = useWorkflowStore(state => state.runWorkflow)
  return useForm({
    defaultValues: defaultWorkflowInput,

    onSubmit: async ({ value }) => {
      const date = new domain.DateRange()

      date.from = value.dateRange.from
      date.to = value.dateRange.to

      const work = new workflow.WorkFlowInput({
        user: await App.GetUser(),
        receive_guides_in_transit: value.receive_guides_in_transit
      })

      work.date = date

      await runWorkflow(
        work,
        new config.BrowserConfig({
          headless: browserForm.state.values.headless,

          trace: browserForm.state.values.trace
        })
      )
    }
  })
}

type AppFormsContextType = {
  browserForm: BrowserConfigForm
  workflowForm: WorkflowForm
}

export const AppFormsContext = createContext<AppFormsContextType | null>(null)

export function useAppForms() {
  const context = useContext(AppFormsContext)

  if (!context) {
    throw new Error('useAppForms must be used inside AppFormsProvider')
  }

  return context
}
