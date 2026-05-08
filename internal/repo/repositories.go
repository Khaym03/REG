package repo

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"
)

type GuideRepository interface {
	Exists(context.Context, domain.DateRange) (bool, error)
	Save(context.Context, domain.DateRange, []domain.Guide) error
	Get(context.Context, domain.DateRange) ([]domain.Guide, error)
}

type ReceptionRepository interface {
	SaveProgress(context.Context, domain.DateRange, domain.ReceptionResult) error
	GetProgress(context.Context, domain.DateRange) (domain.ReceptionResult, error)

	MarkCompleted(context.Context, domain.DateRange) error
	IsCompleted(context.Context, domain.DateRange) (bool, error)
}

type RubroRepository interface {
	Save(context.Context, []domain.Rubro) error
	GetAll(context.Context) ([]domain.Rubro, error)
}
