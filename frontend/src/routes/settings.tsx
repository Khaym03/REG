import { createFileRoute } from '@tanstack/react-router'
import SettiongSection from '@/features/workflow/settings'

export const Route = createFileRoute('/settings')({
  component: SettiongSection
})
