import { useState, type PropsWithChildren } from 'react'
import {
  AppFormsContext,
  useBrowserConfigFormInstance,
  useWorkflowFormInstance
} from '../hooks/use-app'
import { useWorkflowTopics } from '../features/workflow/hooks/use-topics'

export function AppFormsProvider({ children }: PropsWithChildren) {
  const browserForm = useBrowserConfigFormInstance()
  const [isWorkflowRunning, setIsWorkflowRunning] = useState<boolean>(false)
  const workflowForm = useWorkflowFormInstance(
    browserForm,
    setIsWorkflowRunning
  )

  const { currentState } = useWorkflowTopics()

  return (
    <AppFormsContext.Provider
      value={{
        browserForm,
        workflowForm,
        currentState,
        isWorkflowRunning,
        setIsWorkflowRunning
      }}
    >
      {children}
    </AppFormsContext.Provider>
  )
}
