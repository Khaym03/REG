import { useEffect, type PropsWithChildren } from 'react'
import {
  AppFormsContext,
  useBrowserConfigFormInstance,
  useWorkflowFormInstance
} from '../hooks/use-app'
import { useWorkflowStore } from '@/features/workflow/store/index'
import { useUserFormInstance } from '@/hooks/user-form'

export function AppFormsProvider({ children }: PropsWithChildren) {
  const browserForm = useBrowserConfigFormInstance()
  const workflowForm = useWorkflowFormInstance(browserForm)
  const userForm = useUserFormInstance()

  const initWorkflow = useWorkflowStore(state => state.initWorkflow)
  const cleanupListeners = useWorkflowStore(state => state.cleanupListeners)

  // Initialize listeners and build nodes/edges layouts on component mount
  useEffect(() => {
    initWorkflow()

    // Automatically remove backend event hooks when leaving this view
    return () => cleanupListeners()
  }, [initWorkflow, cleanupListeners])


  return (
    <AppFormsContext.Provider
      value={{
        browserForm,
        workflowForm,
        userForm
      }}
    >
      {children}
    </AppFormsContext.Provider>
  )
}
