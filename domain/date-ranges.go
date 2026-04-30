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

func (d DateRange) Key() string {
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

func MonthlyDateRanges(from, to time.Time) []DateRange {
	var ranges []DateRange

	// Normalize to first day of month
	current := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, from.Location())

	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, to.Location())

	for !current.After(end) {
		firstDay := current
		lastDay := current.AddDate(0, 1, -1)

		ranges = append(ranges, DateRange{
			From: firstDay,
			To:   lastDay,
		})

		current = current.AddDate(0, 1, 0)
	}

	return ranges
}
