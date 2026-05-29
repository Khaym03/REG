import { useEffect, type PropsWithChildren } from 'react'
import {
  AppFormsContext,
  useBrowserConfigFormInstance,
  useWorkflowFormInstance
} from '../hooks/use-app'
import { useWorkflowStore } from '@/features/workflow/store/index'

export function AppFormsProvider({ children }: PropsWithChildren) {
  const browserForm = useBrowserConfigFormInstance()
  const workflowForm = useWorkflowFormInstance(browserForm)

  const initWorkflow = useWorkflowStore(state => state.initWorkflow)
  const cleanupListeners = useWorkflowStore(state => state.cleanupListeners)

  // Initialize listeners and build nodes/edges layouts on component mount
  useEffect(() => {
    initWorkflow()

    // Automatically remove backend event hooks when leaving this view
    return () => cleanupListeners()
  }, [initWorkflow, cleanupListeners])

  // const currentState = useWorkflowStore(state => state.currentState)
  // const initListeners = useWorkflowStore(state => state.initListeners)
  // const cleanupListeners = useWorkflowStore(state => state.cleanupListeners)

  // useEffect(() => {
  //   // Initialize the Wails event listeners when the app boots up
  //   initListeners()

  //   // Teardown cleanly when the app closes/hot-reloads
  //   return () => cleanupListeners()
  // }, [initListeners, cleanupListeners])

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
