package e2e

import (
	"context"
	"testing"

	"github.com/Khaym03/REG/app/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
	"github.com/Khaym03/REG/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	testutil.RodSuite
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (suite *LoginTestSuite) TestLoginSuccess() {
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set")
	}

	scraper := &scraper.LoginScraper{
		Browser: suite.Browser,
	}

	handler := command.NewLoginHandler(scraper)

	cmd := command.LoginCommand{
		User: domain.User{
			Username: username,
			Password: password,
		},
	}

	err := handler.Handle(context.Background(), cmd)

	require.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	scraper := &scraper.LoginScraper{
		Browser: suite.Browser,
	}

	handler := command.NewLoginHandler(scraper)

	cmd := command.LoginCommand{
		User: domain.User{
			Username: "fake@example.com",
			Password: "fakepassword",
		},
	}

	err := handler.Handle(context.Background(), cmd)

	assert.Error(suite.T(), err)
}
