package app

import (
	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	"github.com/mustafaturan/bus/v3"
)

type Application struct {
	SessionProvider *auth.Provider
	SessionManger   auth.SessionManager

	EventBus *bus.Bus
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
