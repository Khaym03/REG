import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { TooltipProvider } from '@/components/ui/tooltip'
import { ThemeProvider } from '@/components/theme-provider'
import { createRouter, RouterProvider } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen.ts'

const router = createRouter({ routeTree })

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <TooltipProvider>
        <RouterProvider router={router} />
      </TooltipProvider>
    </ThemeProvider>
  </StrictMode>
)
