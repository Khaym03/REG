package container

import (
	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
)

func buildInventoryHandler(
	repo domain.GuideRepository,
) dcommand.CommandHandler[command.SyncInventoryCommand] {

	base := command.NewInventoryHandler(
		repo,
		scraper.NewInventoryScraper(),
	)

	return withLogging(base)
}
