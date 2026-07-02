package e2e

import (
	"os"
	"path/filepath"

	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type RodSuite struct {
	suite.Suite
	SessionMediator mediator.SessionMediator
}

func (suite *RodSuite) SetupSuite() {
	envPath, _ := filepath.Abs("../.env")
	err := godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}

	// Optional execution guard for network-dependent integration tests
	if os.Getenv("REG_E2E") != "1" {
		suite.T().Skip("Skipping rod-based tests; set REG_E2E=1 to run")
	}
}

func (suite *RodSuite) SetupTest() {
	suite.SessionMediator = mediator.NewSessionMediator(
		event.NewBus(),
	)
}

func (suite *RodSuite) TearDownTest() {
	_ = suite.SessionMediator.Close()
}

func (suite *RodSuite) LoadCredential() (string, string) {
	return os.Getenv("REG_TEST_USERNAME"), os.Getenv("REG_TEST_PASSWORD")
}

func (suite *RodSuite) NewBrowser() *rod.Browser {
	datadir := "rod_data"

	suite.T().TempDir()
	path := filepath.Join(suite.T().TempDir(), datadir)

	l := launcher.New().
		Headless(os.Getenv("REG_HEADLESS") == "1").
		Devtools(false).
		Leakless(false).
		UserDataDir(path)

	return rod.New().
		Context(suite.T().Context()).
		ControlURL(l.MustLaunch()).
		Trace(os.Getenv("REG_ROD_VERBOSE") == "1").
		MustConnect()
}
