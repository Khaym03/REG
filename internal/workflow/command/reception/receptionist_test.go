package reception

import (
	"context"
	"testing"
	"time"

	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/event"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubReceptionRepo struct {
	completed map[string]bool
	saved     []domain.DateRange
	marked    []domain.DateRange
}

func (s *stubReceptionRepo) IsCompleted(_ context.Context, d domain.DateRange) (bool, error) {
	return s.completed[d.From.Format("2006-01")], nil
}

func (s *stubReceptionRepo) SaveProgress(_ context.Context, d domain.DateRange, _ domain.ReceptionResult) error {
	s.saved = append(s.saved, d)
	return nil
}

func (s *stubReceptionRepo) GetProgress(_ context.Context, _ domain.DateRange) (domain.ReceptionResult, error) {
	return domain.ReceptionResult{}, nil
}

func (s *stubReceptionRepo) MarkCompleted(_ context.Context, d domain.DateRange) error {
	s.marked = append(s.marked, d)
	return nil
}

type stubReceptionService struct {
	received []ReceptionOptions
}

func (s *stubReceptionService) Receive(_ context.Context, _ Session, opt ReceptionOptions) (ReceptionResult, error) {
	s.received = append(s.received, opt)
	return ReceptionResult{Processed: 1, Completed: true}, nil
}

func TestReceptionistHandlerPassesMonthlyDateRange(t *testing.T) {
	repo := &stubReceptionRepo{
		completed: map[string]bool{
			"2026-02": true,
		},
	}
	svc := &stubReceptionService{}
	handler := NewReceptionistHandler(
		repo,
		svc,
		logrus.NewEntry(logrus.StandardLogger()),
		event.NewFakeBus(),
	)

	from := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.Local)
	to := time.Date(2026, time.March, 20, 0, 0, 0, 0, time.Local)

	err := handler.Handle(context.Background(), nil, ReceptionistCommand{
		Date: domain.DateRange{
			From: from,
			To:   to,
		},
		ReceiveGuidesInTransit: true,
	})
	require.NoError(t, err)

	// Feb is already completed and not the current month, so Jan and Mar should be processed.
	require.Len(t, svc.received, 2)

	// Each call must receive its single-month DateRange (d), not the full multi-month cmd.Date.
	assert.Equal(t, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.Local), svc.received[0].Date.From)
	assert.Equal(t, time.Date(2026, time.January, 31, 0, 0, 0, 0, time.Local), svc.received[0].Date.To)
	assert.True(t, svc.received[0].ReceiveGuidesInTransit)

	assert.Equal(t, time.Date(2026, time.March, 1, 0, 0, 0, 0, time.Local), svc.received[1].Date.From)
	assert.Equal(t, time.Date(2026, time.March, 31, 0, 0, 0, 0, time.Local), svc.received[1].Date.To)
	assert.True(t, svc.received[1].ReceiveGuidesInTransit)
}
