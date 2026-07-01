import { create } from 'zustand'
import type { RootStoreState } from './types'
import { createFlowSlice } from './flow-slice'
import { createWorkflowSlice } from './workflow-slice'

export const useWorkflowStore = create<RootStoreState>()((...a) => ({
  ...createFlowSlice(...a),
  ...createWorkflowSlice(...a)
}))
