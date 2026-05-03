package e2e

import (
	"testing"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
	"github.com/Khaym03/REG/scraper/session"
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
	provider := session.NewProvider(suite.NewBrowser(), scraper.NewLoginScraper())
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set")
	}

	user := domain.User{
		Username: username,
		Password: password,
	}

	s, err := provider.Start(suite.T().Context(), user)
	require.NoError(suite.T(), err)
	defer func() {
		require.NoError(suite.T(), s.Close())
	}()

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), s)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	provider := session.NewProvider(suite.NewBrowser(), scraper.NewLoginScraper())
	user := domain.User{
		Username: "fake@example.com",
		Password: "wrong",
	}

	s, err := provider.Start(suite.T().Context(), user)
	defer func() {
		require.NoError(suite.T(), s.Close())
	}()

	require.Error(suite.T(), err)
	require.Nil(suite.T(), s)
}
