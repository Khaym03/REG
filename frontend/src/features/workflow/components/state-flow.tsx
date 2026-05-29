import { useEffect } from 'react'
import { ReactFlow, Background } from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { useTheme } from '../../../hooks/use-theme'
// import DevTools from './devtools'
import { LoadingNode } from './node/loading-node'
import { SuccessNode } from './node/success-node'
import { InitialNode } from './node/initial-node'
import { useWorkflowStore } from '@/features/workflow/store/index'

const nodeTypes = {
  initialNode: InitialNode,
  loadingNode: LoadingNode,
  successNode: SuccessNode
}

const ZOOM_FACTOR = 0.8

export default function StateFlow() {
  const { theme } = useTheme()
  const isWorkflowRunning = useWorkflowStore(state => state.isWorkflowRunning)

  // Core React Flow UI State & Handlers (from FlowSlice)
  const nodes = useWorkflowStore(state => state.nodes)
  const edges = useWorkflowStore(state => state.edges)
  const onNodesChange = useWorkflowStore(state => state.onNodesChange)
  const onEdgesChange = useWorkflowStore(state => state.onEdgesChange)
  const onConnect = useWorkflowStore(state => state.onConnect)
  const updateNodeTypes = useWorkflowStore(state => state.updateNodeTypes)

  // Lifecycle & Wails Active Tracking States (from WorkflowSlice)
  const currentState = useWorkflowStore(state => state.currentState)
  const stateHistory = useWorkflowStore(state => state.stateHistory)

  useEffect(() => {
    updateNodeTypes(isWorkflowRunning, currentState, stateHistory)
  }, [isWorkflowRunning, currentState, stateHistory, updateNodeTypes])

  return (
    <>
      <ReactFlow
        colorMode={theme}
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        fitView
        minZoom={ZOOM_FACTOR}
        maxZoom={ZOOM_FACTOR}
        panOnDrag={false}
        zoomOnScroll={false}
        zoomOnPinch={false}
        zoomOnDoubleClick={false}
        elementsSelectable={false}
      >
        <Background color="" />
        {/* <DevTools /> */}
      </ReactFlow>
    </>
  )
}
