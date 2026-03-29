package e2e

import (
	"context"
	"testing"

	"github.com/Khaym03/REG/app/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	RodSuite
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (suite *LoginTestSuite) TestLoginSuccess() {
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set")
	}

	scraper := scraper.NewLoginScraper()

	loginHandler := command.NewLoginHandler(scraper)
	logoutHandler := command.NewLogoutHandler(scraper)

	loginCmd := command.LoginCommand{
		User: domain.User{
			Username: username,
			Password: password,
		},
		Page: suite.Page,
	}

	err := loginHandler.Handle(suite.T().Context(), loginCmd)
	require.NoError(suite.T(), err)

	err = logoutHandler.Handle(suite.T().Context(), command.LogoutCommand{
		Page: suite.Page,
	})
	require.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	scraper := scraper.NewLoginScraper()

	handler := command.NewLoginHandler(scraper)

	cmd := command.LoginCommand{
		User: domain.User{
			Username: "fake@example.com",
			Password: "fakepassword",
		},
		Page: suite.Page,
	}

	err := handler.Handle(context.Background(), cmd)

	assert.Error(suite.T(), err)
}
