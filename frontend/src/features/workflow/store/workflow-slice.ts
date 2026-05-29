import type { StateCreator } from 'zustand'
import { Position, type Node, type Edge } from '@xyflow/react'
import { RunWorkflow, StopWorkflow, Topics } from 'wails/go/main/App'
import { EventsOn } from 'wails/runtime/runtime'
import type { RootStoreState, WorkflowSlice } from './types'
import type { config, workflow } from 'wails/go/models'
import { runDebouncer, stopDebouncer } from '@/lib/utils'

const X_ORIGIN = -400
const X_GAP = 200
const Y_GAP = 75
const nodeDefaults = {
  sourcePosition: Position.Right,
  targetPosition: Position.Left
}

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

  runWorkflow: async (
    work: workflow.WorkFlowInput,
    conf: config.BrowserConfig
  ) => {
    if (get().isWorkflowRunning || get().isDebouncing) return

    await runDebouncer(
      async () => await RunWorkflow(work, conf),
      1000,
      isWaiting => set({ isDebouncing: isWaiting })
    )

    set({ isDebouncing: false })
  },

  stopWorkflow: async () => {
    if (!get().isWorkflowRunning) return

    await stopDebouncer(
      async () => await StopWorkflow(),
      1000,
      isWaiting => set({ isDebouncing: isWaiting })
    )

    set({ isDebouncing: false })
  },

  initWorkflow: async () => {
    get().cleanupListeners()
    const topics = await Topics()

    const initialNodes: Node[] = [
      {
        id: topics.building_browser,
        position: { x: X_ORIGIN, y: 0 },
        data: { label: topics.building_browser },
        sourcePosition: Position.Bottom,
        type: 'initialNode'
      },
      {
        id: topics.login,
        position: { x: X_ORIGIN, y: Y_GAP },
        data: { label: topics.login },
        sourcePosition: Position.Right,
        targetPosition: Position.Top,
        type: 'initialNode'
      },
      {
        id: topics.guides_gather,
        position: { x: X_GAP + X_ORIGIN, y: Y_GAP },
        data: { label: topics.guides_gather },
        type: 'initialNode',
        ...nodeDefaults
      },
      {
        id: topics.inventory_sync,
        position: { x: X_GAP * 2 + X_ORIGIN, y: Y_GAP },
        data: { label: topics.inventory_sync },
        type: 'initialNode',
        ...nodeDefaults
      },
      {
        id: topics.reception,
        position: { x: X_GAP * 3 + X_ORIGIN, y: Y_GAP },
        data: { label: topics.reception },
        type: 'initialNode',
        ...nodeDefaults
      },
      {
        id: topics.logout,
        position: { x: X_GAP * 4 + X_ORIGIN, y: Y_GAP },
        data: { label: topics.logout },
        sourcePosition: Position.Bottom,
        targetPosition: Position.Left,
        type: 'initialNode'
      },
      {
        id: topics.destroying_browser,
        position: { x: X_GAP * 4 + X_ORIGIN, y: Y_GAP * 2 },
        data: { label: topics.destroying_browser },
        targetPosition: Position.Top,
        type: 'initialNode'
      }
    ]

    const initialEdges: Edge[] = [
      {
        id: `${topics.building_browser}-${topics.login}`,
        source: topics.building_browser,
        target: topics.login
      },
      {
        id: `${topics.login}-${topics.guides_gather}`,
        source: topics.login,
        target: topics.guides_gather
      },
      {
        id: `${topics.guides_gather}-${topics.inventory_sync}`,
        source: topics.guides_gather,
        target: topics.inventory_sync
      },
      {
        id: `${topics.inventory_sync}-${topics.reception}`,
        source: topics.inventory_sync,
        target: topics.reception
      },
      {
        id: `${topics.reception}-${topics.logout}`,
        source: topics.reception,
        target: topics.logout
      },
      {
        id: `${topics.logout}-${topics.destroying_browser}`,
        source: topics.logout,
        target: topics.destroying_browser
      }
    ]

    // Initialize cross-slice workflow nodes
    get().setElements(initialNodes, initialEdges)
    set({ stateHistory: [], currentState: '' })

    // Update types instantly using current slice values
    get().updateNodeTypes(get().isWorkflowRunning, '', [])

    const activeTopics = [
      topics.workflow_started,
      topics.building_browser,
      topics.login,
      topics.guides_gather,
      topics.inventory_sync,
      topics.reception,
      topics.logout,
      topics.destroying_browser,
      topics.workflow_finished
    ]

    const localUnsubscribers: (() => void)[] = []

    activeTopics.forEach(event => {
      const unsubscribe = EventsOn(event, () => {
        console.log('[Backend] fire event: ', event)
        const nextHistory = [...get().stateHistory, event]

        set({ currentState: event, stateHistory: nextHistory })

        if (event == topics.workflow_started) {
          set({ isWorkflowRunning: true })
        }

        if (event == topics.workflow_finished) {
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
