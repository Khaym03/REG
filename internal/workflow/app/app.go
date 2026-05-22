package app

import (
	"github.com/Khaym03/REG/app/commands/stats"
	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
)

type Application struct {
	SessionProvider *auth.Provider
	Commands        Commands
	Queries         Queries
}

type Commands struct {
	GatherGuides  guide.GatherGuidesHandler
	SyncInventory inventory.SyncInventoryHandler
	Receptionist  reception.ReceptionistHandler
}

type Queries struct {
	Stats stats.StatsHandler
}
