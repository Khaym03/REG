import { useEffect, useState } from "react"

import type { RegisterUsers } from "bindings/github.com/Khaym03/REG/internal/auth"

import { useAuthStore } from "@/auth/auth-store"
import { useAppForms } from "@/hooks/use-app"
import { cn } from "@/lib/utils"

import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { GetUserPassword } from "bindings/github.com/Khaym03/REG/accountsapi"
import RecentAccounts from "./recent-accounts"

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const { userForm } = useAppForms()

  const getRegisterUsers = useAuthStore((state) => state.getRegisterUsers)

  const [registerUsers, setRegisterUsers] = useState<RegisterUsers[]>([])

  useEffect(() => {
    getRegisterUsers().then(setRegisterUsers)
  }, [getRegisterUsers])

  const handleSelect = async (username: string) => {
    try {
      const user = await GetUserPassword(username)
      userForm.setFieldValue("username", user.username)
      userForm.setFieldValue("password", user.password)
    } catch (e) {
      console.log(e)
    }
  }

  return (
    <div
      className={cn("flex w-2xl flex-col gap-6", className)}
      {...props}
    >
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:grid-cols-2">
          <form
            className="p-6 md:p-8"
            onSubmit={(e) => {
              e.preventDefault()
              e.stopPropagation()
              userForm.handleSubmit()
            }}
          >
            <FieldGroup>
              {/* Header */}
              <div className="flex flex-col items-center gap-2 text-center">
                <h1 className="text-2xl font-bold tracking-tight">Welcome back</h1>
                <p className="text-sm text-muted-foreground text-balance">
                  Login to your account
                </p>
              </div>


              {/* Username */}
              <userForm.Field name="username">
                {(field) => {
                  const isInvalid =
                    field.state.meta.isTouched &&
                    !field.state.meta.isValid

                  return (
                    <Field>
                      <FieldLabel htmlFor={field.name}>
                        Email
                      </FieldLabel>

                      <Input
                        id={field.name}
                        name={field.name}
                        type="text"
                        placeholder="m@example.com"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) =>
                          field.handleChange(e.target.value)
                        }
                        aria-invalid={isInvalid}
                        required
                      />

                      {isInvalid && (
                        <FieldError
                          errors={field.state.meta.errors}
                        />
                      )}
                    </Field>
                  )
                }}
              </userForm.Field>

              {/* Password */}
              <userForm.Field name="password">
                {(field) => {
                  const isInvalid =
                    field.state.meta.isTouched &&
                    !field.state.meta.isValid

                  return (
                    <Field>
                      <FieldLabel htmlFor={field.name}>
                        Password
                      </FieldLabel>

                      <Input
                        id={field.name}
                        name={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) =>
                          field.handleChange(e.target.value)
                        }
                        aria-invalid={isInvalid}
                        required
                      />

                      {isInvalid && (
                        <FieldError
                          errors={field.state.meta.errors}
                        />
                      )}
                    </Field>
                  )
                }}
              </userForm.Field>

              {/* Submit */}
              <userForm.Subscribe
                selector={(state) => [
                  state.canSubmit,
                  state.isSubmitting,
                ]}
              >
                {([_, isSubmitting]) => (
                  <Button
                    type="submit"
                    disabled={isSubmitting}
                    className="transition-all"
                  >
                    {isSubmitting ? (
                      <span className="shiny inline-block bg-[linear-gradient(120deg,rgba(255,255,255,0)_40%,rgba(255,255,255,0.8)_50%,rgba(255,255,255,0)_60%)] bg-size-[200%_100%] bg-clip-text text-white/70 dark:bg-[linear-gradient(120deg,rgba(0,0,0,0)_40%,rgba(0,0,0,0.8)_50%,rgba(0,0,0,0)_60%)] dark:text-foreground/60">
                        Checkeando...
                      </span>
                    ) : (
                      "Submit"
                    )}
                  </Button>
                )}
              </userForm.Subscribe>
            </FieldGroup>
          </form>


          <RecentAccounts users={registerUsers} handleSelect={handleSelect} />
        </CardContent>
      </Card>
    </div>
  )
}
