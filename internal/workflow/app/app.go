package app

import (
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
)

type Application struct {
	SessionMediator mediator.SessionMediator

	EventBus event.Bus
	Commands Commands
	Queries  Queries
}

type Commands struct {
	GatherGuides  guide.GatherGuidesHandler
	SyncInventory inventory.SyncInventoryHandler
	Receptionist  reception.ReceptionistHandler
}

type Queries struct {
	Stats stats.StatsHandler
}
