import { useEffect, useState } from 'react'

import type { RegisterUsers } from 'bindings/github.com/Khaym03/REG/internal/auth'

import { useAuthStore } from '@/auth/auth-store'
import { cn } from '@/lib/utils'

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { FieldGroup } from '@/components/ui/field'
import { GetUserPassword } from 'bindings/github.com/Khaym03/REG/accountsapi'
import RecentAccounts from './recent-accounts'
import { LoadingButton } from './ui/loading-button'

import { router } from '@/routes/router'
import { useAppForm } from '@/hooks/form'
import { User } from 'bindings/github.com/Khaym03/REG/internal/auth'
import * as z from 'zod'

const LoginFormSchema = z.object({
  username: z.string().min(5),
  password: z.string().min(8)
})

const defaultUser = new User({ username: '', password: '' })

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<'div'>) {
  const { navigate } = router
  const login = useAuthStore(s => s.login)
  const isAuthenticated = useAuthStore(s => s.isAuthenticated)

  const form = useAppForm({
    defaultValues: defaultUser,
    onSubmit: async ({ value }) => {
      await login(value.username, value.password)

      if (isAuthenticated()) {
        navigate({ to: '/' })
      }
    },
    validators: {
      onChange: LoginFormSchema,
      onBlur: LoginFormSchema
    }
  })

  const { AppField, setFieldValue, handleSubmit, Subscribe } = form

  const getRegisterUsers = useAuthStore(state => state.getRegisterUsers)

  const [registerUsers, setRegisterUsers] = useState<RegisterUsers[]>([])

  useEffect(() => {
    getRegisterUsers().then(setRegisterUsers)
  }, [getRegisterUsers])

  const handleSelect = async (username: string) => {
    form.reset()

    try {
      const user = await GetUserPassword(username)
      setFieldValue('username', user.username)
      setFieldValue('password', user.password)
    } catch (e) {
      console.log(e)
    }
  }

  return (
    <div className={cn('flex w-2xl flex-col gap-6', className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <form
            className="p-6 md:p-8"
            onSubmit={e => {
              e.preventDefault()
              e.stopPropagation()
              handleSubmit()
            }}
          >
            <FieldGroup>
              <CardHeader
                className="flex flex-col items-center gap-2 text-center"
              >
                <CardTitle className="text-2xl font-bold tracking-tight">
                  Welcome back
                </CardTitle>
                <CardDescription className="text-sm text-muted-foreground text-balance">
                  Login to your account
                </CardDescription>
              </CardHeader>

              <AppField name="username">
                {({ TextField }) => (
                  <TextField label="user" placeholder="j12345" />
                )}
              </AppField>

              <AppField name="password">
                {({ TextField }) => (
                  <TextField
                    label="Password"
                    placeholder="12345"
                    type="password"
                  />
                )}
              </AppField>

              <Subscribe
                selector={state => [state.canSubmit, state.isSubmitting]}
              >
                {([, isSubmitting]) => (
                  <LoadingButton
                    type="submit"
                    isLoading={isSubmitting}
                    disabled={isSubmitting}
                    loadingType="button"
                    loadingText="Checkeando..."
                  >
                    Submit
                  </LoadingButton>
                )}
              </Subscribe>
            </FieldGroup>
          </form>

          <RecentAccounts users={registerUsers} handleSelect={handleSelect} />
        </CardContent>
      </Card>
    </div>
  )
}
