import { format } from 'date-fns/format'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldDescription
} from '@/components/ui/field'
import type { DateRange } from 'bindings/github.com/Khaym03/REG/internal/domain'

export default function DisplaySelectedDate({ dates }: { dates: DateRange }) {
  const startDateStr = format(dates.from, 'MMM yyyy')
  const endDateStr = format(dates.to, 'MMM yyyy')
  return (
    <FieldLabel>
      <Field>
        <FieldContent>
          Selected date:
          <FieldDescription>
            {startDateStr} - {endDateStr}
          </FieldDescription>
        </FieldContent>
      </Field>
    </FieldLabel>
  )
}
