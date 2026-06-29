import { LoginForm } from '@/components/login-form'
import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    // Redirect if already authenticated
    const { user } = context.auth.getState()

    if (user?.logged) {
      throw redirect({ to: '/' })
    }
  },
  component: LoginForm
})
