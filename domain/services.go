package domain

import "context"

type GuideCollector interface {
	Collect(ctx context.Context, dr DateRange) ([]Guide, error)
}

type RubroExtractor interface {
	FromGuides(ctx context.Context, guides []Guide) ([]Rubro, error)
}

type ReceptionService interface {
	Receive(ctx context.Context, dr DateRange) (ReceptionResult, error)
}

type InventoryService interface {
	Snapshot(ctx context.Context) ([]Rubro, error)
	Insert(ctx context.Context, rubro Rubro) error
}
