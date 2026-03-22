package login

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Khaym03/REG/testutil"
)

type LoginTestSuite struct {
	testutil.RodSuite
}

func (suite *LoginTestSuite) TestLoginSuccess() {
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set in .env file")
	}

	user := User{Username: username, Password: password}

	closeSession, err := Login(suite.Page, user)

	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), closeSession)

	if closeSession != nil {
		closeSession()
	}
	time.Sleep(time.Second)
}

func (suite *LoginTestSuite) TestLoginFailureFakeUser() {
	user := User{Username: "fake@example.com", Password: "fakepassword"}
	closeSession, err := Login(suite.Page, user)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), closeSession)
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
