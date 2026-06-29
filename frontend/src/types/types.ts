import type { User } from 'bindings/github.com/Khaym03/REG/internal/auth'
import type { DateRange } from 'bindings/github.com/Khaym03/REG/internal/domain'

export interface WorkflowInput {
  dateRange: DateRange
  receive_guides_in_transit: boolean
}

export interface NavLogoData {
  name: string
  logo: React.ReactNode
  plan: string
}

export interface NavItem {
  name: string
  url: string
  icon: React.ReactNode
}

export interface NavData {
  user: User
  data: NavLogoData
  navItems: NavItem[]
}
