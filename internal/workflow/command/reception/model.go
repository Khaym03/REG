package reception

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/session"
)

type ReceptionOptions struct {
	Date                   DateRange
	ReceiveGuidesInTransit bool
}

type ReceptionService interface {
	Receive(context.Context, Session, ReceptionOptions) (ReceptionResult, error)
}

type Page interface {
	Open() error
	ApplyFilters(DateRange) error
	Rows() ([]TableRow, error)
	ConfirmReception() error
}

// TableRow represents a single line in the results table
type TableRow interface {
	ID() (string, error)
	IsExpired() bool
	TriggerReception() error
}

type (
	Session         = session.Session
	DateRange       = domain.DateRange
	ReceptionResult = domain.ReceptionResult
)
