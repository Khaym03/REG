package e2e

import (
	"testing"
	"time"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/session"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const profileSelector = `#profileDropdown`

type LoginTestSuite struct {
	RodSuite
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (suite *LoginTestSuite) TestLoginSuccess() {
	t := suite.T()
	user := suite.validUser()

	s := suite.mustLogin(user)
	defer func() {
		require.NoError(t, s.Close())
	}()

	err := s.Do(t.Context(), func(p *rod.Page) error {
		_, err := p.Timeout(10 * time.Second).Element(profileSelector)
		return err
	})

	require.NoError(t, err)
}

func (suite *LoginTestSuite) TestSessionAlreadyOpen() {
	t := suite.T()
	user := suite.validUser()

	first := suite.mustLogin(user)
	defer func() {
		_ = suite.SessionMediator.Logout(t.Context(), first)
		require.NoError(t, first.Close())
	}()

	browser := suite.NewBrowser()
	defer browser.MustClose()

	page := browser.MustPage()
	second := session.NewRodSession(page, event.NewBus())
	defer func() {
		require.NoError(t, second.Close())
	}()

	err := suite.SessionMediator.Login(t.Context(), second, user)
	require.ErrorIs(t, err, auth.ErrSessionIsOpen)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	t := suite.T()

	s, err := suite.SessionMediator.Create(t.Context())
	require.NoError(t, err)

	defer func() {
		require.NoError(t, s.Close())
	}()

	err = suite.SessionMediator.Login(t.Context(), s, auth.User{
		Username: "fake@example.com",
		Password: "wrong",
	})

	require.ErrorIs(t, err, auth.ErrInvalidCrendentials)
}

func (suite *LoginTestSuite) mustLogin(user auth.User) session.Session {
	t := suite.T()
	t.Helper()

	s, err := suite.SessionMediator.Create(t.Context())
	require.NoError(t, err)

	err = suite.SessionMediator.Login(t.Context(), s, user)
	require.NoError(t, err)

	return s
}

func (suite *LoginTestSuite) validUser() auth.User {
	t := suite.T()

	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		t.Skip("credentials missing")
	}

	return auth.User{
		Username: username,
		Password: password,
	}
}
