package guide

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/session"
)

type GuideCollector interface {
	Collect(context.Context, Session, domain.DateRange) ([]domain.Guide, error)
}

type RubroExtractor interface {
	FromGuides(context.Context, Session, []Guide) ([]Rubro, error)
}

type (
	Session        = session.Session
	SessionFactory = session.SessionFactory
	DateRange      = domain.DateRange
	Guide          = domain.Guide
	Rubro          = domain.Rubro
)
