import { toast } from 'sonner'

const duration = 1000 * 10

export default function notifyErrToUI(err: unknown) {
  if (!(err instanceof Error)) return

  const e = JSON.parse(err.message) as Error
  const message = e.message.split(':').at(-1)
  toast.error(message, { position: 'bottom-right', duration })
}
