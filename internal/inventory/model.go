package inventory

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/domain"
)

type InventoryService interface {
	Snapshot(context.Context, Session) ([]Rubro, error)
	Insert(context.Context, Session, Rubro) error
}

type (
	Session = auth.Session
	Rubro   = domain.Rubro
)
