package container

import (
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
)

func buildAuthService() domain.AuthService {
	return scraper.NewLoginScraper()
}
