import { useAuthStore } from '@/auth/auth-store'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar
} from '@/components/ui/sidebar'
import { router } from '@/routes/router'
import type { User } from 'bindings/github.com/Khaym03/REG/internal/auth'
import { CaretUpDownIcon, SignOutIcon, UserIcon } from '@phosphor-icons/react'

export function NavUser({ user }: { user: User }) {
  const { isMobile } = useSidebar()
  const logout = useAuthStore(s => s.logout)
  const { navigate } = router

  const handleLogout = async () => {
    await logout()
    navigate({ to: '/login' })
  }
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <DisplayUser user={user} />
              <CaretUpDownIcon className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? 'bottom' : 'right'}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <DisplayUser user={user} />
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout}>
              <SignOutIcon />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

function DisplayUser({ user }: { user: User }) {
  return (
    <>
      <Avatar className="h-8 w-8 rounded-lg">
        <AvatarFallback className="rounded-lg">
          <UserIcon />
        </AvatarFallback>
      </Avatar>
      <div className="grid flex-1 text-left text-sm leading-tight">
        <span className="truncate font-medium">{censore(user.username)}</span>
        <input
          type="password"
          value={user.password}
          className="truncate text-xs"
          readOnly
        />
      </div>
    </>
  )
}

function censore(name: string): string {
  const word = name.split('')
  if (word.length <= 3) return name

  return word.slice(0, 3).join('') + '*'.repeat(word.length - 3)
}
