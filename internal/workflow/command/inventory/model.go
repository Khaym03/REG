package inventory

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/session"
)

type InventoryService interface {
	Snapshot(context.Context, Session) ([]Rubro, error)
	Insert(context.Context, Session, Rubro) error
}

type (
	Session = session.Session
	Rubro   = domain.Rubro
)
