import { createRootRouteWithContext, Outlet } from '@tanstack/react-router'
import type { CSSProperties } from 'react'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { AppSidebar } from '@/components/layout/app-sidebar'

import { useAuthStore } from '@/auth/auth-store'

import { AppFormsProvider } from '@/providers/app-provider'
import TitleBarActions from '@/components/layout/title-bar-actions'

export interface RouterContext {
  auth: typeof useAuthStore
}

// createRootRoute defines the top-level layout
export const Route = createRootRouteWithContext<RouterContext>()({
  component: function RootLayout() {
    return (
      <AppFormsProvider>
        <div className="overflow-hidden">
          <SidebarProvider>
            <SidebarInset>
              <div
                style={{ '--wails-draggable': 'drag' } as CSSProperties}
                className="draggable h-8 w-full flex justify-end fixed top-0 right-0 z-1000"
              >
                <TitleBarActions />
              </div>
              <div className="relative flex justify-center items-center flex-1 flex-col p-4 pt-0 overflow-y-hidden mt-8  gap-4">
                {/* Outlet renders the matching child route */}
                <Outlet />
              </div>
            </SidebarInset>
            <AppSidebar side="right" />
          </SidebarProvider>
        </div>
      </AppFormsProvider>
    )
  }
})
