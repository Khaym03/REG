import { ModeToggle } from './mode-toggle'
import { Button } from '@/components/ui/button'
import { MinusIcon, ResizeIcon, XIcon } from '@phosphor-icons/react'
import { Application, Window } from '@wailsio/runtime'

export default function TitleBarActions() {
  return (
    <div className="grid grid-cols-4 w-[12rem]">
      <ModeToggle />
      <Button
        variant={'ghost'}
        className="border-0 hover:bg-accent"
        onClick={() => Window.Minimise()}
      >
        <MinusIcon />
      </Button>
      <Button
        variant={'ghost'}
        className="border-0 hover:bg-accent"
        onClick={() => Window.ToggleMaximise()}
      >
        <ResizeIcon />
      </Button>
      <Button
        variant={'ghost'}
        className="border-0 hover:bg-destructive"
        onClick={() => Application.Quit()}
      >
        <XIcon />
      </Button>
    </div>
  )
}
