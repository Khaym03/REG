package domain

import "context"

type GuideRepository interface {
	Exists(context.Context, DateRange) (bool, error)
	Save(context.Context, DateRange, []Guide) error
	Get(context.Context, DateRange) ([]Guide, error)
}

type ReceptionRepository interface {
	SaveProgress(context.Context, DateRange, ReceptionResult) error
	GetProgress(context.Context, DateRange) (ReceptionResult, error)

	MarkCompleted(context.Context, DateRange) error
	IsCompleted(context.Context, DateRange) (bool, error)
}

type RubroRepository interface {
	Save(context.Context, []Rubro) error
	GetAll(context.Context) ([]Rubro, error)
}
