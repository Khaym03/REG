import { useAuthStore } from '@/auth/auth-store'
import type { NavData } from '@/types/types'
import { GearIcon, RowsIcon, TreeStructureIcon } from '@phosphor-icons/react'
import { useEffect, useState } from 'react'

const defaults: NavData = {
  data: {
    name: 'REG',
    logo: <RowsIcon />,
    plan: 'free'
  }
}

export default function useNavData(): { navData: NavData } {
  const isAuthenticated = useAuthStore(s => s.user?.logged)
  const user = useAuthStore(s => s.user)

  const [navData, setNavData] = useState<NavData>(defaults)

  useEffect(() => {
    if (isAuthenticated && user) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setNavData(prev => ({
        ...prev,
        user: {
          username: user.username,
          password: ''
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
      }))
    } else {
      setNavData(defaults)
    }
  }, [user, isAuthenticated])

  return {
    navData
  }
}
