package container

import (
	"github.com/Khaym03/REG/internal/auth"
)

func buildAuthService() auth.AuthService {
	return auth.NewLoginScraper()
}
