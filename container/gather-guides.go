package container

import (
	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
)

func buildGatherGuidesHandler(
	guideRepo domain.GuideRepository,
	rubroRepo domain.RubroRepository,
) dcommand.CommandHandler[command.GatherGuidesCommand] {

	scraperSvc := scraper.NewGuidesScraper()
	worker := scraper.NewRodRubroWorker(1)

	base := command.NewGatherGuidesHandler(
		guideRepo,
		rubroRepo,
		scraperSvc,
		worker,
	)

	return withLogging(base)
}
