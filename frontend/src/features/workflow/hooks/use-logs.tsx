import { useEffect } from 'react'
import { useWorkflowStore } from '../store'
import { Events } from '@wailsio/runtime'

export default function useLogs() {
  const addLogLine = useWorkflowStore(state => state.addLogLine)

  useEffect(() => {
    let unsubscribe: (() => void) | undefined

    const setupListener = async () => {
      try {
        unsubscribe = Events.On('LOG', event => {
          const line = event.data as string

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
