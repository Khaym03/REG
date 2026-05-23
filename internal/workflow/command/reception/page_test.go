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
			name: "Successfully processes the first valid row",
			opt: reception.ReceptionOptions{
				Date:                   mockDate,
				ReceiveGuidesInTransit: true,
			},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(nil)
				mp.EXPECT().ApplyFilters(mockDate).Return(nil)

				mp.EXPECT().Rows().Return([]reception.TableRow{mr}, nil)

				// The row executes the trigger, and the page confirms.
				mr.EXPECT().TriggerReception().Return(nil)
				mp.EXPECT().ConfirmReception().Return(nil)
			},
			expectedBool:  true,
			expectedCount: 1,
			expectedErr:   "",
		},
		{
			name: "Skips row when ReceiveGuidesInTransit is false",
			opt: reception.ReceptionOptions{
				Date:                   mockDate,
				ReceiveGuidesInTransit: false,
			},
			setupMocks: func(mp *reception.MockPage, mr *reception.MockTableRow) {
				mp.EXPECT().Open().Return(nil)
				mp.EXPECT().ApplyFilters(mockDate).Return(nil)
				mp.EXPECT().Rows().Return([]reception.TableRow{mr}, nil)
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
