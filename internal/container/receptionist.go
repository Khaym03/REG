package container

import (
	dcommand "github.com/Khaym03/REG/internal/common/decorator/command"
	"github.com/Khaym03/REG/internal/reception"
	"github.com/Khaym03/REG/internal/repo"
)

func buildReceptionHandler(
	repo repo.ReceptionRepository,
) dcommand.CommandHandler[reception.ReceptionistCommand] {
	base := reception.NewReceptionistHandler(
		repo,
		reception.NewReceptionistScraper(),
	)

	return withLogging(base)
}
