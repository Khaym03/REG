import type { PropsWithChildren } from 'react'
import {
  AppFormsContext,
  useBrowserConfigFormInstance,
  useWorkflowFormInstance
} from './use-app'

export function AppFormsProvider({ children }: PropsWithChildren) {
  const browserForm = useBrowserConfigFormInstance()

  const workflowForm = useWorkflowFormInstance(browserForm)

  return (
    <AppFormsContext.Provider
      value={{
        browserForm,
        workflowForm
      }}
    >
      {children}
    </AppFormsContext.Provider>
  )
}
