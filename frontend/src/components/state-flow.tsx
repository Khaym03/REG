import { useState, useCallback, useEffect } from 'react'
import {
  ReactFlow,
  applyNodeChanges,
  applyEdgeChanges,
  addEdge,
  Background,
  Position,
  type Node,
  type Edge,
  type OnConnect,
  type OnEdgesChange,
  type OnNodesChange
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Topics } from 'wails/go/main/App'
import { useTheme } from './use-theme'
// import DevTools from './devtools'
import { LoadingNode } from './loading-node'
import { useWorkflowTopics } from './use-topics'
import { useAppForms } from './use-app'
import { SuccessNode } from './success-node'
import { InitialNode } from './initial-node'

const nodeTypes = {
  initialNode: InitialNode,
  loadingNode: LoadingNode,
  successNode: SuccessNode
}

const X_ORIGIN = -400
const ZOOM_FACTOR = 0.8
const X_GAP = 200
const Y_GAP = 75

const nodeDefaults = {
  sourcePosition: Position.Right,
  targetPosition: Position.Left
}

const getNodes = async (): Promise<Node[]> => {
  const topics = await Topics()

  return [
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
}

const getEdges = async (): Promise<Edge[]> => {
  const topics = await Topics()
  return [
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
}

export default function StateFlow() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [edges, setEdges] = useState<Edge[]>([])

  const { theme } = useTheme()

  useEffect(() => {
    const update = async () => {
      setNodes(await getNodes())
      setEdges(await getEdges())
    }

    update()
  }, [])

  const onNodesChange: OnNodesChange = useCallback(
    changes =>
      setNodes(nodesSnapshot => applyNodeChanges(changes, nodesSnapshot)),
    []
  )
  const onEdgesChange: OnEdgesChange = useCallback(
    changes =>
      setEdges(edgesSnapshot => applyEdgeChanges(changes, edgesSnapshot)),
    []
  )
  const onConnect: OnConnect = useCallback(
    params => setEdges(edgesSnapshot => addEdge(params, edgesSnapshot)),
    []
  )
  const { currentState, stateHistory } = useWorkflowTopics()
  const { isWorkflowRunning } = useAppForms()
  useEffect(() => {
    if (!isWorkflowRunning) {
      setNodes(prevNodes =>
        prevNodes.map(node => ({ ...node, type: 'initialNode' }))
      )
      return
    }

    setNodes(prevNodes =>
      prevNodes.map(node => {
        if (currentState === node.id) {
          return { ...node, type: 'loadingNode' }
        }

        if (stateHistory.includes(node.id)) {
          return { ...node, type: 'successNode' }
        }
        return { ...node, type: 'initialNode' }
      })
    )
  }, [currentState, isWorkflowRunning])

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
