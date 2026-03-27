package domain

import (
	"context"

	"github.com/Khaym03/REG/utils"
)

type Rubro struct {
	Name string `json:"name"`
}

type GuideScraper interface {
	CollectIDs(ctx context.Context, date utils.DateRange) ([]string, error)
}
type RubroWorker interface {
	Process(ctx context.Context, ids []string) ([]Rubro, error)
}

type GuideRepository interface {
	Exists(date utils.DateRange) bool
	SaveIDs(date utils.DateRange, ids []string)
	GetIDs(date utils.DateRange) []string

	SaveRubros([]Rubro)
	GetRubros() []Rubro
}
