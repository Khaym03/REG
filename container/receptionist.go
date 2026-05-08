package container

import (
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/internal/reception"
)

func buildReceptionHandler(
	repo domain.ReceptionRepository,
) dcommand.CommandHandler[reception.ReceptionistCommand] {
	base := reception.NewReceptionistHandler(
		repo,
		reception.NewReceptionistScraper(),
	)

	return withLogging(base)
}
