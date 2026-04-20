package container

import (
	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
)

func buildReceptionHandler(
	repo domain.GuideRepository,
) dcommand.CommandHandler[command.ReceptionistCommand] {

	base := command.NewReceptionistHandler(
		repo,
		scraper.NewReceptionistScraper(),
	)

	return withLogging(base)
}
