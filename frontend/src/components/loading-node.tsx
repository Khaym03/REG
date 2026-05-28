import { BaseNode, BaseNodeContent } from '@/components/base-node'
import { NodeStatusIndicator } from '@/components/node-status-indicator'
import { Handle, type Node, type NodeProps } from '@xyflow/react'

type LodingNode = Node<{ label: string }, 'label'>

export const LoadingNode = ({
  data,
  sourcePosition,
  targetPosition
}: NodeProps<LodingNode>) => {
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
