package utils

import "time"

type DateRange struct {
	From time.Time
	To   time.Time
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
