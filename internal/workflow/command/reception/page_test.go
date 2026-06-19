package reception_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Khaym03/REG/internal/workflow/command/reception"
)

func TestProcessNextExpiredGuide(t *testing.T) {
	mockDate := reception.DateRange{}

	tests := []struct {
		name          string
		opt           reception.ReceptionOptions
		setupMocks    func(mp *reception.MockPage, mr *reception.MockTableRow)
		expectedBool  bool
		expectedCount uint8
		expectedErr   string
	}{
		{
			name: "Successfully processes an expired row even if ReceiveGuidesInTransit is false",
			opt: reception.ReceptionOptions{
				Date:                   mockDate,
				ReceiveGuidesInTransit: false, // Flag is false, but row IS expired
			},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(nil)
				mp.EXPECT().ApplyFilters(mockDate).Return(nil)
				mp.EXPECT().Rows().Return([]reception.TableRow{mr}, nil)

				mr.EXPECT().IsExpired().Return(true) // Always process expired rows
				mr.EXPECT().TriggerReception().Return(nil)
				mp.EXPECT().ConfirmReception().Return(nil)
			},
			expectedBool:  true,
			expectedCount: 1,
			expectedErr:   "",
		},
		{
			name: "Successfully processes an in-transit row when ReceiveGuidesInTransit is true",
			opt: reception.ReceptionOptions{
				Date:                   mockDate,
				ReceiveGuidesInTransit: true, // Flag is true, row is NOT expired
			},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(nil)
				mp.EXPECT().ApplyFilters(mockDate).Return(nil)
				mp.EXPECT().Rows().Return([]reception.TableRow{mr}, nil)

				mr.EXPECT().IsExpired().Return(false)
				mr.EXPECT().TriggerReception().Return(nil)
				mp.EXPECT().ConfirmReception().Return(nil)
			},
			expectedBool:  true,
			expectedCount: 1,
			expectedErr:   "",
		},
		{
			name: "Skips row when NOT expired and ReceiveGuidesInTransit is false",
			opt: reception.ReceptionOptions{
				Date:                   mockDate,
				ReceiveGuidesInTransit: false,
			},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(nil)
				mp.EXPECT().ApplyFilters(mockDate).Return(nil)
				mp.EXPECT().Rows().Return([]reception.TableRow{mr}, nil)

				mr.EXPECT().IsExpired().Return(false) // Not expired + flag false = skip
			},
			expectedBool:  false,
			expectedCount: 0,
			expectedErr:   "",
		},
		{
			name: "Fails when Open returns an error",
			opt:  reception.ReceptionOptions{Date: mockDate},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(errors.New("navigation failed"))
			},
			expectedBool:  false,
			expectedCount: 0,
			expectedErr:   "navigation failed",
		},
		{
			name: "Fails when TriggerReception errors out",
			opt: reception.ReceptionOptions{
				Date:                   mockDate,
				ReceiveGuidesInTransit: true,
			},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(nil)
				mp.EXPECT().ApplyFilters(mockDate).Return(nil)
				mp.EXPECT().Rows().Return([]reception.TableRow{mr}, nil)

				mr.EXPECT().IsExpired().Return(false)
				mr.EXPECT().TriggerReception().Return(errors.New("click failed"))
			},
			expectedBool:  false,
			expectedCount: 0,
			expectedErr:   "trigger reception: click failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPage := reception.NewMockPage(t)
			mockRow := reception.NewMockTableRow(t)

			result := &reception.ReceptionResult{Processed: 0}

			tt.setupMocks(mockPage, mockRow)

			processed, err := reception.ProcessNextExpiredGuide(mockPage, tt.opt, result)

			assert.Equal(t, tt.expectedBool, processed)
			assert.Equal(t, tt.expectedCount, result.Processed)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
