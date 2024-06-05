package handlers

import (
	"os"
	"testing"

	"github.com/ajaen4/go-standard-lib-api/internal/db"
	"github.com/ajaen4/go-standard-lib-api/pkg/handlers"
)

func setupTestCfg(t *testing.T) *handlers.ApiConfig {
	t.Helper()
	testDB, err := db.NewDB("./test_database.json")
	if err != nil {
		t.Fatalf("Error initializing test DB: %s", err)
	}

	apiCfg := &handlers.ApiConfig{
		DB:        testDB,
		JwtSecret: os.Getenv("JWT_SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}

	return apiCfg
}

func tearDownTestCfg(t *testing.T, apiCfg *handlers.ApiConfig) {
	t.Helper()
	apiCfg.DB.RemoveDB()
}
