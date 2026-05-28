import { BaseNode, BaseNodeContent } from '@/components/base-node'
import { NodeStatusIndicator } from '@/components/node-status-indicator'
import { Handle, type Node, type NodeProps } from '@xyflow/react'

type successNode = Node<{ label: string }, 'label'>

export const SuccessNode = ({
  data,
  sourcePosition,
  targetPosition
}: NodeProps<successNode>) => {
  return (
    <NodeStatusIndicator status="success" variant="border">
      <BaseNode className=" ">
        <BaseNodeContent>{data.label}</BaseNodeContent>
      </BaseNode>

      {targetPosition && <Handle type="target" position={targetPosition} />}

      {sourcePosition && <Handle type="source" position={sourcePosition} />}
    </NodeStatusIndicator>
  )
}
