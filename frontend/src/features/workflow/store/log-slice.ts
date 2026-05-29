import type { StateCreator } from 'zustand'
import type { RootStoreState, LogSlice, LogEntry } from './types'

export const createLogSlice: StateCreator<
  RootStoreState,
  [],
  [],
  LogSlice
> = set => ({
  entries: [],

  addLogLine: (line, maxEntries = 500) =>
    set(state => {
      let newEntry: LogEntry

      try {
        const data = JSON.parse(line)
        newEntry = {
          id: crypto.randomUUID(),
          time: new Date().toLocaleTimeString(),
          level: data.level?.toLowerCase() || 'info',
          message: data.msg || data.message || line
        }
      } catch {
        newEntry = {
          id: crypto.randomUUID(),
          time: new Date().toLocaleTimeString(),
          level: 'info',
          message: line
        }
      }

      const updatedEntries = [...state.entries, newEntry]
      return {
        entries: updatedEntries.slice(-maxEntries)
      }
    }),

  clearLogs: () => set({ entries: [] })
})
