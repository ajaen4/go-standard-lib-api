package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ajaen4/go-standard-lib-api/internal/db"
	"github.com/ajaen4/go-standard-lib-api/pkg/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := db.NewDB("./database.json")
	if err != nil {
		log.Fatalf("Error initializing DB: %s", err)
	}

	apiCfg := &handlers.ApiConfig{
		DB:        db,
		JwtSecret: os.Getenv("JWT_SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}

	mux := http.NewServeMux()
	handlers.AssignHandlers(mux, apiCfg)
}
