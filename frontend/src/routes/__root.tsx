import { createRootRouteWithContext, Outlet } from '@tanstack/react-router'
import type { CSSProperties } from 'react'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { AppSidebar } from '@/components/layout/app-sidebar'

import { useAuthStore } from '@/auth/auth-store'

import { AppFormsProvider } from '@/providers/app-provider'

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
                className="draggable h-7 w-full bg-background flex justify-end fixed top-0 right-0 z-0"
              ></div>
              <div className="relative flex justify-center items-center flex-1 flex-col p-4 py-0 overflow-y-hidden mt-7 h-[572px] max-h-[572px] gap-4">
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
