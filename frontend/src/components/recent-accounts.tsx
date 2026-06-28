import { cn } from "@/lib/utils"
import { CaretRightIcon, ClockIcon, UserCircleIcon } from "@phosphor-icons/react"
import type { RegisterUsers } from "bindings/github.com/Khaym03/REG/internal/auth"

interface RecentAccountsProps {
  users: RegisterUsers[]
  handleSelect: (username: string) => Promise<void>
}

export default function RecentAccounts({ users, handleSelect }: RecentAccountsProps) {

  return (

    <aside className="hidden bg-muted/50 border-l border-border md:flex flex-col">
      <div className="px-6 pt-6 pb-3">
        <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
          Recent accounts
        </p>
      </div>

      <div className="flex flex-col gap-1 px-3 pb-6">
        {users.map((user) => (
          <button
            key={user.username}
            type="button"
            onClick={() => handleSelect(user.username)}
            className={cn(
              "group flex w-full items-center gap-3 rounded-lg px-3 py-3",
              "text-left transition-colors",
              "hover:bg-accent hover:text-accent-foreground",
              "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            )}
          >
            <span className="flex size-9 shrink-0 items-center justify-center rounded-full bg-background border border-border text-muted-foreground group-hover:border-primary/30 group-hover:text-primary transition-colors">
              <UserCircleIcon className="size-4" />
            </span>

            <span className="flex min-w-0 flex-1 flex-col gap-0.5">
              <span className="truncate text-sm font-medium leading-none">
                {user.username}
              </span>
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                {user.last_use ? (
                  <>
                    <ClockIcon className="size-3" />
                    Last used
                  </>
                ) : (
                  "Never used"
                )}
              </span>
            </span>

            <CaretRightIcon className="size-4 shrink-0 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
          </button>
        ))}
      </div>
    </aside>
  )
}
