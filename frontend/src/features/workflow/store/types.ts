import type {
  Node,
  Edge,
  OnNodesChange,
  OnEdgesChange,
  OnConnect
} from '@xyflow/react'
import * as config from 'bindings/github.com/Khaym03/REG/internal/config'
import * as workflow from 'bindings/github.com/Khaym03/REG/internal/workflow'

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

export interface LogEntry {
  id: string
  time: string
  level: string
  message: string
}

export interface LogSlice {
  entries: LogEntry[]
  addLogLine: (line: string, maxEntries?: number) => void
  clearLogs: () => void
}

export type RootStoreState = FlowSlice & WorkflowSlice & LogSlice
