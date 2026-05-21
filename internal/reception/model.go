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

type (
	Session         = auth.Session
	DateRange       = domain.DateRange
	ReceptionResult = domain.ReceptionResult
)
