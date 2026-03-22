package testutil

import (
	"os"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type RodSuite struct {
	suite.Suite
	Browser *rod.Browser
	Page    *rod.Page
}

func (suite *RodSuite) SetupSuite() {
	// Load .env file if it exists; optional and safe to ignore error
	_ = godotenv.Load()

	// Optional execution guard for network-dependent integration tests
	if os.Getenv("REG_E2E") != "1" {
		suite.T().Skip("Skipping rod-based tests; set REG_E2E=1 to run")
	}

	headless := true
	if os.Getenv("REG_HEADLESS") == "0" || os.Getenv("REG_HEADLESS") == "false" {
		headless = false
	}

	l := launcher.New().
		Headless(headless).
		Devtools(false).
		Leakless(false)

	browser := rod.New().
		ControlURL(l.MustLaunch()).
		MustConnect()

	suite.Browser = browser
}

func (suite *RodSuite) TearDownSuite() {
	if suite.Browser != nil {
		suite.Browser.MustClose()
	}
}

func (suite *RodSuite) SetupTest() {
	// Open a fresh page for each test.
	if suite.Browser != nil {
		suite.Page = suite.Browser.MustPage()
		suite.Page.MustNavigate(c.LoginURL)
	}
}

func (suite *RodSuite) TearDownTest() {
	if suite.Page != nil {
		suite.Page.MustClose()
		suite.Page = nil
	}
}

func (suite *RodSuite) LoadCredential() (string, string) {
	return os.Getenv("REG_TEST_USERNAME"), os.Getenv("REG_TEST_PASSWORD")
}
