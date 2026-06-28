import { useAuthStore } from "@/auth/auth-store"
import { router } from "@/routes/router"
import { useForm } from "@tanstack/react-form"
import { User } from "bindings/github.com/Khaym03/REG/internal/auth"

export type UserForm = ReturnType<typeof useUserFormInstance>
const defaultUser = new User({username:"", password: ""})
export function useUserFormInstance() {
  const {navigate} = router
  const login = useAuthStore((s) => s.login)
  const isAuthenticated = useAuthStore(s => s.isAuthenticated)
  return useForm({
    defaultValues: defaultUser,
    onSubmit: async ({ value }) => {
      await login(value.username, value.password)

      if (isAuthenticated()) {
        navigate({to: "/"})
      }
    }
  })
}
