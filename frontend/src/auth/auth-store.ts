import { create } from "zustand"
import { RegisterUsers, User } from "bindings/github.com/Khaym03/REG/internal/auth"
import {AccountsAPI} from "bindings/github.com/Khaym03/REG"

export type AuthState = {
  user: RegisterUsers | null
  registerUsers: RegisterUsers[]
  loading: boolean

  setUser: (user: RegisterUsers | null) => void
  setLoading: (loading: boolean) => void

  getRegisterUsers: () => Promise<RegisterUsers[]>

  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  initialize: () => Promise<void>

  isAuthenticated: () => boolean
}

export const useAuthStore = create<AuthState>((set,get) => ({
  user: null,
  registerUsers: [],
  loading: true,

  setUser: (user) => set({ user }),
  setLoading: (loading) => set({ loading }),

  async getRegisterUsers() {
    let registerUsers: RegisterUsers[] = []
    try {
      registerUsers = await AccountsAPI.GetRegisterUsers()
      set({registerUsers})
    } catch (e) {
      console.log(e)
      set({registerUsers})
    }

    return registerUsers
  },

  async initialize() {
    const user = await AccountsAPI.CurrentUser()
    set({user})
  },

  async login(username, password) {
    try {
     const user  = new User({username,password})
     await AccountsAPI.AuthUser(user)

    set({
        user: await AccountsAPI.CurrentUser(),
      })
    } catch (e) {
      console.log(e)
      set({user: null})
    }
  },

  async logout() {
   const {user} = get()
    if (!user) return

    user.logged = false

    try {
      await AccountsAPI.UpdateUser(user)
    } catch (e) {
      console.log(e)
    }

    set({user: null})
  },

  isAuthenticated() {
    return get().user !== null
  }
}))
