import { createRootRouteWithContext, Outlet, redirect } from '@tanstack/react-router'
import type { CSSProperties } from 'react'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { AppSidebar } from '@/components/layout/app-sidebar'
import { App } from 'bindings/github.com/Khaym03/REG'
import { GearIcon, RowsIcon, TreeStructureIcon } from '@phosphor-icons/react'
import type { NavData } from '@/types/types'
import type { useAuthStore } from '@/auth/auth-store'

import { AppFormsProvider } from '@/providers/app-provider'

export interface RouterContext {
  auth: typeof useAuthStore
}

// createRootRoute defines the top-level layout
export const Route = createRootRouteWithContext<RouterContext>()({
  beforeLoad: async ({ context }) => {
    const auth = context.auth.getState()

    await auth.initialize()

    if (!auth.isAuthenticated()) {
      redirect({
        to: '/login',
        search: {
          redirect: "/",
        },
      })
    } else {
      redirect({ to: "/" })
    }

  },
  loader: getNavData,
  component: function RootLayout() {
    const navData = Route.useLoaderData()
    return (
      <AppFormsProvider>
        <div className="overflow-hidden">
          <SidebarProvider>
            <SidebarInset>
              <div
                style={{ '--wails-draggable': 'drag' } as CSSProperties}
                className="draggable h-7 w-full bg-background flex justify-end fixed top-0 right-0 z-0"
              ></div>
              <div className="relative flex justify-center items-center flex-1 flex-col p-4 py-0 overflow-y-hidden mt-7 h-[572px] max-h-[572px] gap-4">
                {/* Outlet renders the matching child route */}
                <Outlet />
              </div>
            </SidebarInset>
            <AppSidebar navData={navData} side="right" />
          </SidebarProvider>
        </div>
      </AppFormsProvider>
    )
  }
})

async function getNavData(): Promise<NavData> {
  const user = await App.GetUser()
  return {
    user: {
      name: user.username,
      password: user.password,
      avatar: ''
    },
    data: {
      name: 'REG',
      logo: <RowsIcon />,
      plan: 'free'
    },
    navItems: [
      {
        name: 'Workflow',
        url: '/',
        icon: <TreeStructureIcon />
      },
      {
        name: 'Settings',
        url: '/settings',
        icon: <GearIcon />
      }
    ]
  }
}
