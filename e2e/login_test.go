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

	provider *auth.Provider
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (suite *LoginTestSuite) TestLoginSuccess() {
	user := suite.validUser()

	s, err := suite.provider.Start(suite.T().Context(), user, suite.NewBrowser())
	defer suite.closeSession(s)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), s)

	err = s.Do(suite.T().Context(), func(p *rod.Page) error {
		_, err := p.Timeout(10 * time.Second).Element(profileSelector)
		return err
	})

	require.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestSessionAlreadyOpen() {
	user := suite.validUser()

	rootBrowser := suite.NewBrowser()
	helperBrowser := rootBrowser.MustIncognito()

	// First login
	first, err := suite.provider.Start(
		suite.T().Context(),
		user,
		rootBrowser,
	)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), first)

	defer func() {
		require.NoError(suite.T(), first.Close(suite.T().Context()))
	}()

	// Second login while first session is still alive
	second, err := suite.provider.Start(
		suite.T().Context(),
		user,
		helperBrowser,
	)

	require.ErrorIs(suite.T(), err, auth.ErrSessionIsOpen)
	require.Nil(suite.T(), second)

}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	user := auth.User{
		Username: "fake@example.com",
		Password: "wrong",
	}

	s, err := suite.provider.Start(suite.T().Context(), user, suite.NewBrowser())
	defer suite.closeSession(s)

	require.ErrorIs(suite.T(), err, auth.ErrInvalidCrendentials)
	require.Nil(suite.T(), s)
}

func (suite *LoginTestSuite) SetupTest() {
	suite.provider = auth.NewProvider(
		auth.NewLoginScraper(),
		event.NewBus(),
	)
}

func (suite *LoginTestSuite) validUser() auth.User {
	username, password := suite.LoadCredential()

	if username == "" || password == "" {
		suite.T().Skip("credentials missing")
	}

	return auth.User{
		Username: username,
		Password: password,
	}
}

func (suite *LoginTestSuite) closeSession(s auth.Session) {
	if s != nil {
		require.NoError(suite.T(), s.Close(suite.T().Context()))
	}
}

const profileSelector = `#profileDropdown`
