package domain

import (
	"context"
	"fmt"

	"github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/utils"
)

type Rubro struct {
	Name string `json:"name"`
}

type Guide struct {
	ID string
}

func (g Guide) URL() string {
	return fmt.Sprintf("%s/%s", constants.GuidesURL, g.ID)
}

type GuideScraper interface {
	CollectGuides(ctx context.Context, date utils.DateRange) ([]Guide, error)
}
type RubroWorker interface {
	Process(ctx context.Context, ids []Guide) ([]Rubro, error)
}

type GuideRepository interface {
	Exists(utils.DateRange) bool
	SaveGuides(utils.DateRange, []Guide)
	GetGuides(utils.DateRange) []Guide

	// Progress
	SaveReceptionProgress(utils.DateRange, ReceptionResult)
	GetReceptionProgress(utils.DateRange) ReceptionResult

	// Final state
	MarkReceptionCompleted(utils.DateRange)
	IsReceptionCompleted(utils.DateRange) bool

	SaveRubros([]Rubro)
	GetRubros() []Rubro
}

type InventoryScraper interface {
	RubrosSnapshot(context.Context) ([]Rubro, error)
	Insert(context.Context, Rubro) error
}

type Receptionist interface {
	// Receive all the expired [Guide] in a given [utils.DateRange]
	Receive(context.Context, utils.DateRange) (ReceptionResult, error)
}

type ReceptionResult struct {
	Processed int
	Completed bool
}
