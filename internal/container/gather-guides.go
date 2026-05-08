package container

import (
	dcommand "github.com/Khaym03/REG/internal/common/decorator/command"
	"github.com/Khaym03/REG/internal/guide"
	"github.com/Khaym03/REG/internal/repo"
)

func buildGatherGuidesHandler(
	guideRepo repo.GuideRepository,
	rubroRepo repo.RubroRepository,
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
