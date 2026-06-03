import { useEffect } from 'react'
import { useWorkflowStore } from '../store'
import { EventsOn } from 'wails/runtime/runtime'

export default function useLogs() {
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

  return { addLogLine }
}
