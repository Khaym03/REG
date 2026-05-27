import type { PropsWithChildren } from 'react'
import {
  AppFormsContext,
  useBrowserConfigFormInstance,
  useWorkflowFormInstance
} from './use-app'
import { useWorkflowTopics } from './use-topics'

export function AppFormsProvider({ children }: PropsWithChildren) {
  const browserForm = useBrowserConfigFormInstance()

  const workflowForm = useWorkflowFormInstance(browserForm)

  const { currentState } = useWorkflowTopics()

  return (
    <AppFormsContext.Provider
      value={{
        browserForm,
        workflowForm,
        currentState
      }}
    >
      {children}
    </AppFormsContext.Provider>
  )
}
