import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function createManagedDebouncer() {
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  return async <T>(
    asyncFn: () => Promise<T>,
    delay: number,
    onStateChange: (isWaiting: boolean) => void
  ): Promise<T | undefined> => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }

    // Signal immediately to the UI that we are in a debouncing/waiting state
    onStateChange(true)

    return new Promise<T | undefined>((resolve, reject) => {
      timeoutId = setTimeout(async () => {
        try {
          // Clear the debounce state right before execution starts
          onStateChange(false)
          const result = await asyncFn()
          resolve(result)
        } catch (error) {
          reject(error)
        } finally {
          timeoutId = null
        }
      }, delay)
    })
  }
}

// Instantiate dedicated debouncers for both actions outside the store hook
// so they preserve their respective timeoutIDs across state updates.
export const runDebouncer = createManagedDebouncer()
export const stopDebouncer = createManagedDebouncer()
