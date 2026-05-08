package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonthlyDateRanges(t *testing.T) {
	// Mocking "now" to May 15, 2026
	now := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		expected []DateRange
	}{
		{
			name: "Past months only (Full months)",
			from: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC),
			expected: []DateRange{
				{
					From: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
				},
				{
					From: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Includes current month (Should cap To at 'now')",
			from: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
			expected: []DateRange{
				{
					From: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					From: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
					To:   now, // Capped here
				},
			},
		},
		{
			name:     "Entirely in the future (Should cap immediately)",
			from:     time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
			expected: nil, // Loop shouldn't even start or ranges stay empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MonthlyDateRanges(tt.from, tt.to, now)

			t.Logf("Running test: %s", tt.name)
			t.Logf("Input range: from %s to %s (Now: %s)",
				tt.from.Format("2006-01-02"),
				tt.to.Format("2006-01-02"),
				now.Format("2006-01-02"),
			)

			for i, r := range result {
				// This uses your DateRange.String() implementation
				t.Logf("Range [%d]: %s", i, r)
			}

			// Assertions using testify
			if tt.expected == nil {
				assert.Empty(t, result)
			} else {
				assert.Len(t, result, len(tt.expected))
				for i := range tt.expected {
					assert.True(t, tt.expected[i].From.Equal(result[i].From), "From date mismatch at index %d", i)
					assert.True(t, tt.expected[i].To.Equal(result[i].To), "To date mismatch at index %d", i)
				}
			}
		})
	}
}
