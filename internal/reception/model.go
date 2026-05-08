package reception

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/domain"
)

type ReceptionService interface {
	Receive(context.Context, Session, DateRange) (ReceptionResult, error)
}

type (
	Session         = auth.Session
	DateRange       = domain.DateRange
	ReceptionResult = domain.ReceptionResult
)
