package e2e

import (
	"testing"
	"time"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/event"
	"github.com/go-rod/rod"
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
	provider := auth.NewProvider(
		auth.NewLoginScraper(),
		event.NewBus(),
	)
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set")
	}

	user := auth.User{
		Username: username,
		Password: password,
	}

	s, err := provider.Start(suite.T().Context(), user, suite.NewBrowser())
	defer func() {
		require.NoError(suite.T(), s.Close(suite.T().Context()))
	}()

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), s)

	err = s.Do(suite.T().Context(), func(p *rod.Page) error {
		_, err := p.Timeout(10 * time.Second).Element(profileSelector)
		return err
	})

	require.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	provider := auth.NewProvider(

		auth.NewLoginScraper(),
		event.NewBus(),
	)
	user := auth.User{
		Username: "fake@example.com",
		Password: "wrong",
	}

	s, err := provider.Start(suite.T().Context(), user, suite.NewBrowser())
	defer func() {
		if s != nil {
			require.NoError(suite.T(), s.Close(suite.T().Context()))
		}
	}()

	require.ErrorIs(suite.T(), err, auth.ErrInvalidCrendentials)
	require.Nil(suite.T(), s)
}

const profileSelector = `#profileDropdown`
