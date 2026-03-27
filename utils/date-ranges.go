package utils

import "time"

type DateRange struct {
	From time.Time
	To   time.Time
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

func MonthlyDateRangesCurrentToLastYear() []DateRange {
	var dateSet []DateRange

	now := time.Now()

	// Loop to generate the set from current month back to last year
	for i := 0; i <= 12; i++ {
		// Calculate the month relative to now
		targetDate := now.AddDate(0, -i, 0)

		// Get the 1st day of that month
		firstDay := time.Date(targetDate.Year(), targetDate.Month(), 1, 0, 0, 0, 0, targetDate.Location())

		// Get the last day (1st day of next month minus 1 day)
		lastDay := firstDay.AddDate(0, 1, -1)

		// Append to our set
		dateSet = append(dateSet, DateRange{
			From: firstDay,
			To:   lastDay,
		})
	}

	return dateSet
}
