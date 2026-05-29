import type {
  Node,
  Edge,
  OnNodesChange,
  OnEdgesChange,
  OnConnect
} from '@xyflow/react'

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
  currentState: string
  stateHistory: string[]
  unsubscribers: (() => void)[]

  initWorkflow: () => Promise<void>
  cleanupListeners: () => void
}

export type RootStoreState = FlowSlice & WorkflowSlice
