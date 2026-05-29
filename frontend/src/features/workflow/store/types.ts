import type {
  Node,
  Edge,
  OnNodesChange,
  OnEdgesChange,
  OnConnect
} from '@xyflow/react'
import type { config, workflow } from 'wails/go/models'

export interface FlowSlice {
  nodes: Node[]
  edges: Edge[]
  onNodesChange: OnNodesChange
  onEdgesChange: OnEdgesChange
  onConnect: OnConnect
  setElements: (nodes: Node[], edges: Edge[]) => void
  updateNodeTypes: (
    isWorkflowRunning: boolean,
    currentState: string,
    stateHistory: string[]
  ) => void
}

export interface WorkflowSlice {
  isWorkflowRunning: boolean
  isDebouncing: boolean
  currentState: string
  stateHistory: string[]
  unsubscribers: (() => void)[]

  runWorkflow: (
    work: workflow.WorkFlowInput,
    conf: config.BrowserConfig
  ) => Promise<void>
  stopWorkflow: () => Promise<void>
  initWorkflow: () => Promise<void>
  cleanupListeners: () => void
}

export type RootStoreState = FlowSlice & WorkflowSlice
