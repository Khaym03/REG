import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem
} from '@/components/ui/sidebar'
import { Link } from '@tanstack/react-router'

export interface NavItem {
  name: string
  url: string
  icon: React.ReactNode
}

export function NavREG({ items }: { items: NavItem[] }) {
  return (
    <SidebarGroup className="group-data-[collapsible=icon]:hidden">
      <SidebarMenu>
        {items.map(item => (
          <SidebarMenuItem key={item.name}>
            <SidebarMenuButton asChild>
              <Link to={item.url}>
                {item.icon}
                <span>{item.name}</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}
