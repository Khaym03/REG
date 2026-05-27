import { useState, useEffect } from 'react'
import { Topics } from 'wails/go/main/App'
import { EventsOn } from 'wails/runtime/runtime'

export const useWorkflowTopics = () => {
  const [currentState, setCurrentState] = useState<string>('')
  const [isLoading, setIsLoading] = useState<boolean>(true)

  useEffect(() => {
    let isMounted = true

    const unsubscribers: (() => void)[] = []

    const setupStatsListener = async () => {
      try {
        setIsLoading(true)

        const topics = await Topics()

        if (!isMounted) return

        const activeTopics = [
          topics.workflow_started_topic,
          topics.loggin_topic,
          topics.guides_gather_topic,
          topics.inventory_sync_topic,
          topics.reception_topic,
          topics.logout_topic,
          topics.workflow_finished_topic
        ]

        console.log('active topics: ', activeTopics)

        activeTopics.forEach(event => {
          unsubscribers.push(
            EventsOn(event, data => {
              console.log(`[Backend Event] Triggered! Topic: ${event}`, data)

              setCurrentState(event)
            })
          )
        })
      } catch (error) {
        console.error('Error configuring backend listeners:', error)
      } finally {
        if (isMounted) {
          setIsLoading(false)
        }
      }
    }

    setupStatsListener()

    return () => {
      isMounted = false

      unsubscribers.forEach(unsubscribe => unsubscribe())
    }
  }, [])

  return { currentState, isLoading }
}
