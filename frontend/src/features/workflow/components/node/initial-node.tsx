import {
  BaseNode,
  BaseNodeContent
} from '@/features/workflow/components/node/base-node'
import { NodeStatusIndicator } from '@/features/workflow/components/node/node-status-indicator'
import { Handle, type Node, type NodeProps } from '@xyflow/react'

type initialNode = Node<{ label: string }, 'label'>

export const InitialNode = ({
  data,
  sourcePosition,
  targetPosition
}: NodeProps<initialNode>) => {
  return (
    <NodeStatusIndicator status="initial" variant="border">
      <BaseNode className=" ">
        <BaseNodeContent>{data.label}</BaseNodeContent>
      </BaseNode>

      {targetPosition && <Handle type="target" position={targetPosition} />}

      {sourcePosition && <Handle type="source" position={sourcePosition} />}
    </NodeStatusIndicator>
  )
}
