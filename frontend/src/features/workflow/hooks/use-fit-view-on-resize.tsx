import { useEffect, useRef } from 'react'
import { useReactFlow } from '@xyflow/react'

export default function FitViewObserver() {
  const { fitView } = useReactFlow()
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const el = ref.current?.parentElement
    if (!el) return

    const observer = new ResizeObserver(() => {
      requestAnimationFrame(() => {
        fitView({
          padding: 0.2,
          duration: 200
        })
      })
    })

    observer.observe(el)

    return () => observer.disconnect()
  }, [fitView])

  return <div ref={ref} />
}
