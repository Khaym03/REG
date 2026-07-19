import type { StateCreator } from 'zustand'
import { type Node, type Edge } from '@xyflow/react'
import { App } from 'bindings/github.com/Khaym03/REG'
import { Events } from '@wailsio/runtime'
import type { RootStoreState, WorkflowSlice } from './types'
import type { BrowserConfig } from 'bindings/github.com/Khaym03/REG/internal/config'
import type { WorkFlowInput } from 'bindings/github.com/Khaym03/REG/internal/workflow'
import { runDebouncer, stopDebouncer } from '@/lib/utils'
import { buildInitialEdges, buildInitialNodes } from '../graph/adapter'
import { Topic } from 'bindings/github.com/Khaym03/REG/internal/event/models'
import notifyErrToUI from '@/lib/notify'

export const createWorkflowSlice: StateCreator<
  RootStoreState,
  [],
  [],
  WorkflowSlice
> = (set, get) => ({
  isWorkflowRunning: false,
  isDebouncing: false,
  currentState: '',
  stateHistory: [],
  unsubscribers: [],

  runWorkflow: async (work: WorkFlowInput, conf: BrowserConfig) => {
    if (get().isWorkflowRunning || get().isDebouncing) return

    const run = async () => {
      try {
        await App.RunWorkflow(work, conf)
      } catch (e) {
        console.error(e)
        notifyErrToUI(e)
      }
    }

    await runDebouncer(run, 1000, isWaiting => set({ isDebouncing: isWaiting }))

    set({ isDebouncing: false })
  },

  stopWorkflow: async () => {
    if (!get().isWorkflowRunning) return

    const stop = async () => {
      try {
        await App.StopWorkflow()
      } catch (e) {
        console.error(e)
        notifyErrToUI(e)
      }
    }
    await stopDebouncer(stop, 1000, isWaiting =>
      set({ isDebouncing: isWaiting })
    )

    set({ isDebouncing: false })
  },

  initWorkflow: async () => {
    get().cleanupListeners()

    const initialNodes: Node[] = buildInitialNodes()
    const initialEdges: Edge[] = buildInitialEdges()

    // Initialize cross-slice workflow nodes
    get().setElements(initialNodes, initialEdges)
    set({ stateHistory: [], currentState: '' })

    // Update types instantly using current slice values
    get().updateNodeTypes(get().isWorkflowRunning, '', [])

    const activeTopics = Object.values(Topic)

    console.assert(activeTopics.length > 7, "Missing Topics")

    const localUnsubscribers: (() => void)[] = []

    activeTopics.forEach(event => {
      const unsubscribe = Events.On(event, () => {
        console.log('[Backend] fire event: ', event)
        const nextHistory = [...get().stateHistory, event]

        set({ currentState: event, stateHistory: nextHistory })

        if (event === Topic.WorkflowStarted) {
          set({ isWorkflowRunning: true })
        }

        if (event === Topic.WorkflowFinished) {
          set({
            isWorkflowRunning: false,
            stateHistory: [],
            isDebouncing: false
          })
        }

        // Notify the flow slice to repaint node styling
        get().updateNodeTypes(get().isWorkflowRunning, event, nextHistory)
      })
      localUnsubscribers.push(unsubscribe)
    })

    set({ unsubscribers: localUnsubscribers })
  },

  cleanupListeners: () => {
    const { unsubscribers } = get()
    if (unsubscribers.length > 0) {
      unsubscribers.forEach(unsub => unsub())
      set({ unsubscribers: [] })
    }
  }
})
