package stats_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Khaym03/REG/internal/stats"
)

func TestStats_HasActionableGuides(t *testing.T) {
	tests := []struct {
		name               string
		stats              stats.Stats
		receiveInTransit   bool
		expectedActionable bool
	}{
		{
			name:               "no guides at all",
			stats:              stats.Stats{},
			receiveInTransit:   false,
			expectedActionable: false,
		},
		{
			name: "expired guides should always be actionable",
			stats: stats.Stats{
				ExpiredGuides: 2,
			},
			receiveInTransit:   false,
			expectedActionable: true,
		},
		{
			name: "in-transit guides enabled",
			stats: stats.Stats{
				InTransitGuides: 3,
			},
			receiveInTransit:   true,
			expectedActionable: true,
		},
		{
			name: "in-transit guides disabled",
			stats: stats.Stats{
				InTransitGuides: 3,
			},
			receiveInTransit:   false,
			expectedActionable: false,
		},
		{
			name: "expired guides take precedence even if in-transit disabled",
			stats: stats.Stats{
				InTransitGuides: 2,
				ExpiredGuides:   1,
			},
			receiveInTransit:   false,
			expectedActionable: true,
		},
		{
			name: "expired and in-transit enabled",
			stats: stats.Stats{
				InTransitGuides: 4,
				ExpiredGuides:   2,
			},
			receiveInTransit:   true,
			expectedActionable: true,
		},
		{
			name: "other stats alone are not actionable",
			stats: stats.Stats{
				OutstandingDebt:   100,
				PendingProcedures: 1,
			},
			receiveInTransit:   true,
			expectedActionable: false,
		},
		{
			name: "all stats populated with in-transit enabled",
			stats: stats.Stats{
				OutstandingDebt:   100,
				InTransitGuides:   5,
				ExpiredGuides:     1,
				PendingProcedures: 2,
			},
			receiveInTransit:   true,
			expectedActionable: true,
		},
		{
			name: "all stats populated with only expired actionable",
			stats: stats.Stats{
				OutstandingDebt:   100,
				InTransitGuides:   5,
				ExpiredGuides:     1,
				PendingProcedures: 2,
			},
			receiveInTransit:   false,
			expectedActionable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.HasActionableGuides(tt.receiveInTransit)

			require.Equal(t, tt.expectedActionable, result)
		})
	}
}
