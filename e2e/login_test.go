package e2e

import (
	"testing"

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

	err := scraper.Login(suite.T().Context(), suite.Page, domain.User{
		Username: username,
		Password: password,
	})
	require.NoError(suite.T(), err)

	err = scraper.Logout(suite.T().Context(), suite.Page)
	require.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	scraper := scraper.NewLoginScraper()

	err := scraper.Login(suite.T().Context(), suite.Page, domain.User{
		Username: "fake@example.com",
		Password: "fakepassword",
	})

	assert.Error(suite.T(), err)
}
