package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/Khaym03/REG/internal/session"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type AccountServiceSuite struct {
	suite.Suite

	service AccountService
}

func (s *AccountServiceSuite) SetupTest() {
	p := repo.NewJSONPersistence(
		filepath.Join(s.T().TempDir(), "test_store_user.json"),
		func() []RegisterUsers { return nil },
	)

	s.service = AccountService{
		p:    p,
		auth: NewLoginScraper(event.NewFakeBus()),
	}

	envPath, _ := filepath.Abs("../../.env")
	err := godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}

}

func (s *AccountServiceSuite) TestAuthUser() {
	if os.Getenv("REG_E2E") != "1" {
		s.T().Skip("Skipping rod-based tests; set REG_E2E=1 to run")
	}

	datadir := "rod_test_data"
	rootDir := filepath.Join(s.T().TempDir(), datadir)

	l := launcher.New().
		Headless(os.Getenv("REG_HEADLESS") == "1").
		Devtools(false).
		Leakless(false).
		UserDataDir(filepath.Join(rootDir, datadir))

	b := rod.New().
		Context(s.T().Context()).
		ControlURL(l.MustLaunch()).
		Trace(os.Getenv("REG_ROD_VERBOSE") == "1").
		MustConnect()

	rodSess := session.NewRodSession(b.MustPage(), event.NewFakeBus())
	defer func() {
		err := rodSess.Close()
		s.Require().NoError(err)
	}()

	user := User{
		Username: os.Getenv("REG_TEST_USERNAME"),
		Password: os.Getenv("REG_TEST_PASSWORD"),
	}
	err := s.service.AuthUser(s.T().Context(), user, rodSess)

	s.Require().NoError(err)
	users, err := s.service.p.Load()
	s.Require().NoError(err)
	s.Assert().Len(len(users), 1)
}

func (s *AccountServiceSuite) TestStoreUserSecret() {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "happy path",
			user: User{
				Username: "test",
				Password: "1234",
			},
		},
		{
			name:    "invalid user",
			user:    User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.service.StoreUserSecret(tt.user)

			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *AccountServiceSuite) TestGetUserPassword() {
	user := User{
		Username: "test",
		Password: "1234",
	}

	err := s.service.StoreUserSecret(user)
	s.Require().NoError(err)

	usr, err := s.service.GetUserPassword(user.Username)
	s.Require().NoError(err)
	s.Assert().Equal(user, usr)
}

func TestAccountServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceSuite))
}
