package domain

import (
	"fmt"
	"time"

	c "github.com/Khaym03/REG/constants"
)

type DateRange struct {
	From time.Time
	To   time.Time
}

func (d DateRange) MonthKey() string {
	return d.From.Format(c.DateKeyFormat)
}

func (d DateRange) String() string {
	const layout = "2006-01-02"
	return fmt.Sprintf(
		"[%s - %s]",
		d.From.Format(layout),
		d.To.Format(layout),
	)
}

func MonthlyDateRanges(from, to, now time.Time) []DateRange {
	var ranges []DateRange

	// If 'to' is in the future, cap it to the current time
	if to.After(now) {
		to = now
	}

	// Normalize to first day of month for consistent iteration
	current := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, from.Location())
	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, to.Location())

	for !current.After(end) {
		firstDay := current
		// Calculate the theoretical last day of the month
		lastDay := current.AddDate(0, 1, -1)

		// Final check: if the calculated last day exceeds 'now',
		// we use 'now' as the limit.
		if lastDay.After(now) {
			lastDay = now
		}

		// Prevent adding a range where From is after To (could happen if 'from' is today)
		if !firstDay.After(lastDay) {
			ranges = append(ranges, DateRange{
				From: firstDay,
				To:   lastDay,
			})
		}

		// Move to the first day of the next month
		current = current.AddDate(0, 1, 0)
	}

	return ranges
}
