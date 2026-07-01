import { Position } from '@xyflow/react'
import { Topic } from 'bindings/github.com/Khaym03/REG/internal/event/models'

export interface WorkflowDefinition {
  label: string
  position: {
    x: number
    y: number
  }
  sourcePosition?: Position
  targetPosition?: Position
}

const nodeDefaults = {
  sourcePosition: Position.Right,
  targetPosition: Position.Left
}

export const WORKFLOW_NODES: Partial<Record<Topic, WorkflowDefinition>> = {
  [Topic.BuildingBrowser]: {
    label: 'Starting browser',
    position: { x: 0, y: 0 },
    sourcePosition: Position.Bottom
  },

  [Topic.Login]: {
    label: 'Logging in',
    position: { x: 0, y: 75 },
    sourcePosition: Position.Right,
    targetPosition: Position.Top
  },

  [Topic.GuidesGather]: {
    label: 'Gathering guides',
    position: { x: 200, y: 75 },
    ...nodeDefaults
  },

  [Topic.InventorySync]: {
    label: 'Syncing inventory',
    position: { x: 400, y: 75 },
    ...nodeDefaults
  },

  [Topic.Reception]: {
    label: 'Receiving guides',
    position: { x: 600, y: 75 },

    ...nodeDefaults
  },

  [Topic.Logout]: {
    label: 'Logging out',
    position: { x: 800, y: 75 },
    sourcePosition: Position.Bottom,
    targetPosition: Position.Left
  },

  [Topic.DestroyingBrowser]: {
    label: 'Closing browser',
    position: { x: 800, y: 150 },
    targetPosition: Position.Top
  }
}

export const WORKFLOW_EDGES: [Topic, Topic][] = [
  [Topic.BuildingBrowser, Topic.Login],
  [Topic.Login, Topic.GuidesGather],
  [Topic.GuidesGather, Topic.InventorySync],
  [Topic.InventorySync, Topic.Reception],
  [Topic.Reception, Topic.Logout],
  [Topic.Logout, Topic.DestroyingBrowser]
]
