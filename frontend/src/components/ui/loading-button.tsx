import { Button, type buttonVariants } from './button'
import type { VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'
import React from 'react'

interface LoadingButtonProps
  extends
    React.ComponentPropsWithoutRef<typeof Button>,
    VariantProps<typeof buttonVariants> {
  isLoading?: boolean
  loadingText?: React.ReactNode
  loadingType?: 'submit' | 'button' | 'reset'
}

const LoadingButton = React.forwardRef<HTMLButtonElement, LoadingButtonProps>(
  (
    {
      isLoading = false,
      loadingText = 'Loading...',
      loadingType,
      disabled,
      children,
      className,
      type = 'button',
      ...props
    },
    ref
  ) => {
    return (
      <Button
        ref={ref}
        type={isLoading && loadingType ? loadingType : type}
        disabled={disabled}
        className={cn('transition-all', className)}
        {...props}
      >
        {isLoading ? (
          <span className="shiny inline-block bg-[linear-gradient(120deg,rgba(255,255,255,0)_40%,rgba(255,255,255,0.8)_50%,rgba(255,255,255,0)_60%)] bg-size-[200%_100%] bg-clip-text text-white/70 dark:bg-[linear-gradient(120deg,rgba(0,0,0,0)_40%,rgba(0,0,0,0.8)_50%,rgba(0,0,0,0)_60%)] dark:text-foreground/60">
            {loadingText}
          </span>
        ) : (
          children
        )}
      </Button>
    )
  }
)

LoadingButton.displayName = 'LoadingButton'

export { LoadingButton }
