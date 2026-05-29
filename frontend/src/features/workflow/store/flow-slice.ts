import type { StateCreator } from 'zustand'
import { applyNodeChanges, applyEdgeChanges, addEdge } from '@xyflow/react'
import type { RootStoreState, FlowSlice } from './types'

export const createFlowSlice: StateCreator<
  RootStoreState,
  [],
  [],
  FlowSlice
> = set => ({
  nodes: [],
  edges: [],

  onNodesChange: changes =>
    set(state => ({
      nodes: applyNodeChanges(changes, state.nodes)
    })),

  onEdgesChange: changes =>
    set(state => ({
      edges: applyEdgeChanges(changes, state.edges)
    })),

  onConnect: params =>
    set(state => ({
      edges: addEdge(params, state.edges)
    })),

  setElements: (nodes, edges) => set({ nodes, edges }),

  updateNodeTypes: (isWorkflowRunning, currentState, stateHistory) => {
    set(state => ({
      nodes: state.nodes.map(node => {
        if (!isWorkflowRunning) {
          return { ...node, type: 'initialNode' }
        }
        if (currentState === node.id) {
          return { ...node, type: 'loadingNode' }
        }
        if (stateHistory.includes(node.id)) {
          return { ...node, type: 'successNode' }
        }
        return { ...node, type: 'initialNode' }
      })
    }))
  }
})
