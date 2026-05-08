package e2e

import (
	"os"
	"path/filepath"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type RodSuite struct {
	suite.Suite
	provider *auth.Provider
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

// func (suite *RodSuite) TearDownSuite() {

// }

// func (suite *RodSuite) SetupTest() {
// 	// Open a fresh page for each test.
// 	if suite.Browser != nil {
// 		suite.Page = suite.Browser.MustPage()
// 		suite.Page.MustNavigate(c.BaseURL)
// 	}
// }

// func (suite *RodSuite) TearDownTest() {
// 	if suite.Page != nil {
// 		suite.Page.MustClose()
// 		suite.Page = nil
// 	}
// }

func (suite *RodSuite) LoadCredential() (string, string) {
	return os.Getenv("REG_TEST_USERNAME"), os.Getenv("REG_TEST_PASSWORD")
}

func (suite *RodSuite) NewBrowser() *rod.Browser {
	envPath, _ := filepath.Abs("../.env")
	rootDir := filepath.Dir(envPath)
	l := launcher.New().
		Headless(os.Getenv("REG_HEADLESS") == "1").
		Devtools(false).
		Leakless(false).
		UserDataDir(filepath.Join(rootDir, "rod_data"))

	return rod.New().
		Context(suite.T().Context()).
		ControlURL(l.MustLaunch()).
		Trace(os.Getenv("REG_ROD_VERBOSE") == "1").
		MustConnect()
}
