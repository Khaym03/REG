import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ThemeProvider } from '@/providers/theme-provider.tsx'
import { RouterProvider } from '@tanstack/react-router'
import { AppFormsProvider } from './providers/app-provider.tsx'
import { useAuthStore } from './auth/auth-store.ts'
import { router } from './routes/router.ts'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <TooltipProvider>
        <AppFormsProvider>
          <RouterProvider
            router={router}
            context={{
              auth: useAuthStore
            }}
          />
        </AppFormsProvider>
      </TooltipProvider>
    </ThemeProvider>
  </StrictMode>
)
