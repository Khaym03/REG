package container

import (
	"github.com/Khaym03/REG/scraper"
	"github.com/Khaym03/REG/session"
)

func buildAuthService() session.AuthService {
	return scraper.NewLoginScraper()
}
