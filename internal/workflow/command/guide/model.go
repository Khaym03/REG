package guide

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/domain"
)

type GuideCollector interface {
	Collect(context.Context, auth.Session, domain.DateRange) ([]domain.Guide, error)
}

type RubroExtractor interface {
	FromGuides(context.Context, Session, []Guide) ([]Rubro, error)
}

type (
	Session   = auth.Session
	DateRange = domain.DateRange
	Guide     = domain.Guide
	Rubro     = domain.Rubro
)
