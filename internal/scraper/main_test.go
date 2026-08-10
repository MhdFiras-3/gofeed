package scraper

import (
	"os"
	"testing"

	"github.com/MhdFiras-3/gofeed/internal/handlers"
	"github.com/MhdFiras-3/gofeed/internal/testingutils"
)

var testCfg *handlers.APIConfig

func TestMain(m *testing.M) {
	queries, _, cleanup := testingutils.SetupTestDB()
	testCfg = &handlers.APIConfig{
		DB: queries,
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}
