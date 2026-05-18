import { NavREG } from '@/components/nav-reg'
import { NavUser } from '@/components/nav-user'
import { NavLogo } from '@/components/nav-logo'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail
} from '@/components/ui/sidebar'
import { XIcon, MinusIcon } from '@phosphor-icons/react'
import { ModeToggle } from './mode-toggle'
import { Button } from './ui/button'
import { Quit, WindowMinimise } from 'wails/runtime/runtime'
import type { ComponentProps } from 'react'
import type { NavData } from '@/types/types'

interface AppSidebarProps extends ComponentProps<typeof Sidebar> {
  navData: NavData
}
export function AppSidebar({ navData, ...props }: AppSidebarProps) {
  return (
    <Sidebar side="right" collapsible="icon" {...props}>
      <div className="grid grid-cols-3">
        <ModeToggle />
        <Button
          variant={'ghost'}
          className="border-0 hover:bg-accent"
          onClick={() => WindowMinimise()}
        >
          <MinusIcon />
        </Button>
        <Button
          variant={'ghost'}
          className="border-0 hover:bg-destructive"
          onClick={() => Quit()}
        >
          <XIcon />
        </Button>
      </div>
      <SidebarHeader>
        {navData && <NavLogo data={navData.data} />}
      </SidebarHeader>
      <SidebarContent>
        {navData && <NavREG items={navData.navItems} />}
      </SidebarContent>
      <SidebarFooter>
        {navData && <NavUser user={navData.user} />}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
