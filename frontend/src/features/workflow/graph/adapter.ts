import type { Edge, Node } from '@xyflow/react'
import { WORKFLOW_EDGES, WORKFLOW_NODES } from './definitions'

export function buildInitialNodes(): Node[] {
  return Object.entries(WORKFLOW_NODES).map(([key, config]) => ({
    id: key,
    type: 'initialNode',
    data: {
      label: config.label
    },
    position: config.position,
    sourcePosition: config.sourcePosition,
    targetPosition: config.targetPosition
  }))
}

export function buildInitialEdges(): Edge[] {
  return WORKFLOW_EDGES.map(([from, to]) => ({
    id: `${from}-${to}`,
    source: from,
    target: to
  }))
}
