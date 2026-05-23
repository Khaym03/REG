package reception

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/domain"
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
	Session         = auth.Session
	DateRange       = domain.DateRange
	ReceptionResult = domain.ReceptionResult
)
