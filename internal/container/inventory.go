package container

import (
	dcommand "github.com/Khaym03/REG/internal/common/decorator/command"
	"github.com/Khaym03/REG/internal/inventory"
	"github.com/Khaym03/REG/internal/repo"
)

func buildInventoryHandler(
	repo repo.RubroRepository,
) dcommand.CommandHandler[inventory.SyncInventoryCommand] {
	base := inventory.NewInventoryHandler(
		repo,
		inventory.NewInventoryScraper(),
	)

	return withLogging(base)
}
