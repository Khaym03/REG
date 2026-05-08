package container

import (
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/internal/guide"
)

func buildGatherGuidesHandler(
	guideRepo domain.GuideRepository,
	rubroRepo domain.RubroRepository,
) dcommand.CommandHandler[guide.GatherGuidesCommand] {
	scraperSvc := guide.NewGuidesScraper()
	worker := guide.NewRodRubroWorker(1)

	base := guide.NewGatherGuidesHandler(
		guideRepo,
		rubroRepo,
		scraperSvc,
		worker,
	)

	return withLogging(base)
}
