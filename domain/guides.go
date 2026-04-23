package domain

import (
	"context"
	"fmt"

	"github.com/Khaym03/REG/constants"
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
	CollectGuides(ctx context.Context, date DateRange) ([]Guide, error)
}
type RubroWorker interface {
	Process(ctx context.Context, ids []Guide) ([]Rubro, error)
}

type GuideRepository interface {
	Exists(DateRange) bool
	SaveGuides(DateRange, []Guide)
	GetGuides(DateRange) []Guide

	// Progress
	SaveReceptionProgress(DateRange, ReceptionResult)
	GetReceptionProgress(DateRange) ReceptionResult

	// Final state
	MarkReceptionCompleted(DateRange)
	IsReceptionCompleted(DateRange) bool

	SaveRubros([]Rubro)
	GetRubros() []Rubro
}

type InventoryScraper interface {
	RubrosSnapshot(context.Context) ([]Rubro, error)
	Insert(context.Context, Rubro) error
}

type Receptionist interface {
	// Receive all the expired [Guide] in a given [utils.DateRange]
	Receive(context.Context, DateRange) (ReceptionResult, error)
}

type ReceptionResult struct {
	Processed int
	Completed bool
}
