import { useEffect, useRef, useState } from 'react'
import { cn } from '@/lib/utils'
import { EventsOn } from 'wails/runtime'
import { useVirtualizer } from '@tanstack/react-virtual'
import { TerminalWindowIcon } from '@phosphor-icons/react'

interface LogEntry {
  id: string
  time: string
  level: string
  message: string
}

const COLORS: Record<string, string> = {
  error: 'text-red-400',
  warning: 'text-amber-400',
  info: 'text-emerald-400',
  debug: 'text-sky-400',
  trace: 'text-zinc-500'
}

export function TerminalLogs({ eventName = 'LOG', maxEntries = 500 }) {
  const [entries, setEntries] = useState<LogEntry[]>([])
  const viewportRef = useRef<HTMLDivElement>(null)

  const rowVirtualizer = useVirtualizer({
    count: entries.length,
    getScrollElement: () => viewportRef.current,
    estimateSize: () => 25,
    useFlushSync: false // Disable synchronous updates
  })

  useEffect(() => {
    // Scroll to bottom whenever entries change
    if (viewportRef.current) {
      const scrollContainer = viewportRef.current.querySelector(
        '[data-radix-scroll-area-viewport]'
      )
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight
      }
    }
  }, [entries])

  useEffect(() => {
    let unsubscribe: (() => void) | undefined

    const setupListener = async () => {
      try {
        unsubscribe = EventsOn(eventName, (line: string) => {
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

          setEntries(prev => [...prev, newEntry].slice(-maxEntries))
        })
      } catch (err) {
        console.error('Failed to attach log listener', err)
      }
    }

    setupListener()
    return () => unsubscribe?.()
  }, [eventName, maxEntries])

  return (
    <>
      <div
        ref={viewportRef}
        data-slot="scroll-area"
        style={{
          height: `266px`,
          width: `799px`,
          overflow: 'auto'
        }}
        className=""
      >
        <div
          style={{
            height: `
            ${rowVirtualizer.getTotalSize() > 0 ? rowVirtualizer.getTotalSize() + 'px' : '100%'}`,
            width: '100%',
            position: 'relative'
          }}
        >
          {entries.length == 0 && (
            <div className="size-full flex justify-center items-center text-foreground/10">
              <TerminalWindowIcon size={96} />
            </div>
          )}
          {rowVirtualizer.getVirtualItems().map(virtualRow => (
            <div
              key={virtualRow.index}
              className={`flex gap-3 leading-relaxed border-b`}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: `${virtualRow.size}px`,
                transform: `translateY(${virtualRow.start}px)`
              }}
            >
              <div className="shrink-0 select-none flex justify-center items-center">
                {entries[virtualRow.index].time}
              </div>
              <div
                className={cn(
                  'uppercase font-bold shrink-0 w-12 flex justify-center items-center',
                  COLORS[entries[virtualRow.index].level] || ''
                )}
              >
                {entries[virtualRow.index].level}
              </div>
              <div className="text-xs leading-none font-medium break-all flex justify-center items-center">
                {entries[virtualRow.index].message}
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}
