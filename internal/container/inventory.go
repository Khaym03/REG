package container

import (
	dcommand "github.com/Khaym03/REG/internal/common/decorator/command"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/inventory"
)

func buildInventoryHandler(
	repo domain.RubroRepository,
) dcommand.CommandHandler[inventory.SyncInventoryCommand] {
	base := inventory.NewInventoryHandler(
		repo,
		inventory.NewInventoryScraper(),
	)

	return withLogging(base)
}
