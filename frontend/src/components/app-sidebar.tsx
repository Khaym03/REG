import * as React from 'react'

import { NavREG, type NavItem } from '@/components/nav-reg'
import { NavUser, type User } from '@/components/nav-user'
import { NavLogo, type NavLogoData } from '@/components/nav-logo'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail
} from '@/components/ui/sidebar'
import { RowsIcon, TreeStructureIcon, GearIcon } from '@phosphor-icons/react'

interface NavData {
  user: User
  data: NavLogoData
  navItems: NavItem[]
}

const data: NavData = {
  user: {
    name: 'shadcn',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg'
  },
  data: {
    name: 'REG',
    logo: <RowsIcon />,
    plan: 'free'
  },
  navItems: [
    {
      name: 'Workflow',
      url: '#',
      icon: <TreeStructureIcon />
    },
    {
      name: 'Settings',
      url: '#',
      icon: <GearIcon />
    }
  ]
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <NavLogo data={data.data} />
      </SidebarHeader>
      <SidebarContent>
        <NavREG items={data.navItems} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
