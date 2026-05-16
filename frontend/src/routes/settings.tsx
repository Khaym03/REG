import { createFileRoute } from '@tanstack/react-router'
import SettiongSection from '@/components/settings'

export const Route = createFileRoute('/settings')({
  component: SettiongSection
})
