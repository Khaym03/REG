import { LoginForm } from '@/components/login-form'
import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    // Redirect if already authenticated
    if (context.auth.getState().isAuthenticated()) {
      throw redirect({ to: "/" })
    }
  },
  component: LoginForm
})
