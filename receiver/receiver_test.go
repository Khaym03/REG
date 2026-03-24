package receiver

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/login"
	"github.com/Khaym03/REG/testutil"
)

type ReceiverTestSuite struct {
	testutil.RodSuite
}

func TestReceiverSuite(t *testing.T) {
	suite.Run(t, new(ReceiverTestSuite))
}

func (suite *ReceiverTestSuite) TestGuidesIDGatherer() {
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set in .env file")
	}

	user := login.User{Username: username, Password: password}

	// Login first
	closeSession, err := login.Login(suite.Page, user)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), closeSession)
	defer closeSession()

	// Navigate to reception page
	suite.Page.MustNavigate(c.ReceptionURL)

	gatherer := &GuidesIDGatherer{
		monthlyGuideIDs: make(monthlyGuideIDs),
	}

	err = gatherer.ApplyFiltersToGuideReceiver(suite.Page)
	assert.NoError(suite.T(), err)

	// Assert that some guides were collected
	assert.NotEmpty(suite.T(), gatherer.monthlyGuideIDs)
}

func (suite *ReceiverTestSuite) TestExtractRubrosFromGuide() {
	absPath, err := filepath.Abs("guide-sample.html")
	require.NoError(suite.T(), err)

	suite.Page.MustNavigate("file://" + absPath)

	rubros := ExtractRubrosFromGuide(suite.Page)

	require.NotEmpty(suite.T(), rubros, "The list of extracted rubros should not be empty")

	expectedFirst := Rubro{
		Nombre:   "ACEITE PARA MANUFACTURA",
		Cantidad: "0,050 TM",
		// PrecioVenta:  "1,00 BS",
		// Presentacion: "OTROS",
		Marca: "OTROS",
	}

	assert.Equal(suite.T(), expectedFirst.Nombre, rubros[0].Nombre)
	assert.Equal(suite.T(), expectedFirst.Cantidad, rubros[0].Cantidad)
	assert.Equal(suite.T(), expectedFirst.Marca, rubros[0].Marca)

	if len(rubros) > 1 {
		assert.Equal(suite.T(), "EDULCORANTES IMPORTADOS", rubros[1].Nombre)
		assert.Equal(suite.T(), "0,002 TM", rubros[1].Cantidad)
	}
}

func (suite *ReceiverTestSuite) TestGuidesIDGathererWithExistingData() {
	username, password := suite.LoadCredential()
	if username == "" || password == "" {
		suite.T().Skip("Skipping test: REG_TEST_USERNAME and REG_TEST_PASSWORD not set in .env file")
	}

	user := login.User{Username: username, Password: password}

	// Login first
	closeSession, err := login.Login(suite.Page, user)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), closeSession)
	defer closeSession()

	// Navigate to reception page
	suite.Page.MustNavigate(c.ReceptionURL)

	gatherer := &GuidesIDGatherer{
		monthlyGuideIDs: monthlyGuideIDs{
			"2026-03": []string{},
			"2026-02": []string{},
			"2026-01": []string{},
			"2025-12": []string{},
			"2025-11": []string{},
			"2025-10": []string{},
			"2025-09": []string{},
			"2025-08": []string{},
			"2025-07": []string{},
			"2025-06": []string{},
			"2025-05": []string{},
			"2025-04": []string{},
			"2025-03": []string{},
		},
	}

	err = gatherer.ApplyFiltersToGuideReceiver(suite.Page)
	assert.NoError(suite.T(), err)

	// Assert that the existing data is preserved (empty map values)
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2026-03"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2026-02"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2026-01"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-12"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-11"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-10"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-09"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-08"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-07"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-06"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-05"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-04"])
	assert.Equal(suite.T(), []string{}, gatherer.monthlyGuideIDs["2025-03"])
	// Assert that data map remains non-empty
	assert.NotEmpty(suite.T(), gatherer.monthlyGuideIDs)
}
