import { createFileRoute, redirect, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: async ({ context }) => {
    const auth = context.auth.getState()

    await auth.initialize()

    const isAuthenticated = auth.user?.logged
    if (!isAuthenticated) {
      throw redirect({
        to: '/login',
        search: {
          redirect: '/'
        }
      })
    }
  },
  component: () => <Outlet />
})
