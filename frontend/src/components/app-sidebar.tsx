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
import {
  XIcon,
  MinusIcon,
  ReceiptIcon,
  BankIcon,
  ReceiptXIcon,
  ClockUserIcon
} from '@phosphor-icons/react'
import { ModeToggle } from './mode-toggle'
import { Button } from './ui/button'
import {
  Quit,
  WindowMinimise,
  EventsOn,
  EventsOff
} from 'wails/runtime/runtime'
import { useEffect, useState, type ComponentProps } from 'react'
import type { NavData } from '@/types/types'
import { Topics } from 'wails/go/main/App'
import { stats } from 'wails/go/models'
import { Label } from './ui/label'

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
        <Stats />
        {navData && <NavUser user={navData.user} />}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}

const statItems = [
  {
    label: 'Deuda Pendiente',
    key: 'outstanding_debt',
    Icon: BankIcon
  },
  {
    label: 'Guías en Tránsito',
    key: 'intransit_guides',
    Icon: ReceiptIcon
  },
  {
    label: 'Guías Vencidas',
    key: 'expired_guides',
    Icon: ReceiptXIcon
  },
  {
    label: 'Trámites Pendientes',
    key: 'pending_procedures',
    Icon: ClockUserIcon
  }
] as const

function Stats() {
  const [currentStats, setCurrentStats] = useState<stats.Stats>(
    new stats.Stats({
      outstanding_debt: 0,
      intransit_guides: 0,
      expired_guides: 0,
      pending_procedures: 0
    })
  )

  useEffect(() => {
    // Variable para controlar si el componente se desmontó mientras esperábamos el async
    let isMounted = true
    let activeTopic = ''

    const setupStatsListener = async () => {
      const topics = await Topics()

      if (topics.stats_result && isMounted) {
        activeTopic = topics.stats_result

        EventsOn(activeTopic, datas => {
          setCurrentStats(new stats.Stats(datas))
        })
      }
    }

    setupStatsListener()

    return () => {
      isMounted = false
      if (activeTopic) {
        EventsOff(activeTopic)
      }
    }
  }, [])
  return (
    <div className="grid grid-rows-4 gap-1">
      {statItems.map(item => (
        <div key={item.key} className="flex gap-4 text-2xl">
          <div className=" flex aspect-square size-8 items-center justify-center rounded-lg text-sidebar-primary-foreground">
            <item.Icon weight="light" />
          </div>
          <div className="flex flex-col flex-3">
            <Label>{item.label}</Label>
            <div className="grid items-center text-base leading-0 py-2.5">
              {currentStats[item.key]}
              {item.key === 'outstanding_debt' && 'Bs'}
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
