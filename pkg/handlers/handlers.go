package handlers

import (
	"log"
	"net/http"

	"github.com/ajaen4/go-standard-lib-api/internal/db"
)

type ApiConfig struct {
	JwtSecret      string
	PolkaKey       string
	FileserverHits int
	DB             *db.DB
}

func AssignHandlers(mux *http.ServeMux, apiCfg *ApiConfig) {
	mux.HandleFunc("GET /api/healthz", HealthCheck)

	fileHandler := http.FileServer(http.Dir("."))

	mux.Handle("GET /app/*", apiCfg.MiddlewareMetricsInc(fileHandler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsCount)
	mux.HandleFunc("/api/reset", apiCfg.MetricsReset)

	mux.HandleFunc("GET /api/chirps", NewHandler(apiCfg.GetChirps))
	mux.HandleFunc("GET /api/chirps/{chirpID}", NewHandler(apiCfg.GetChirp))
	mux.HandleFunc("POST /api/chirps", NewHandler(apiCfg.PostChirp))
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", NewHandler(apiCfg.DeleteChirp))

	mux.HandleFunc("POST /api/users", NewHandler(apiCfg.PostUser))
	mux.HandleFunc("PUT /api/users", NewHandler(apiCfg.PutUser))

	mux.HandleFunc("POST /api/login", NewHandler(apiCfg.PostLogin))
	mux.HandleFunc("POST /api/refresh", NewHandler(apiCfg.PostRefToken))
	mux.HandleFunc("POST /api/revoke", NewHandler(apiCfg.PostRevokeToken))

	mux.HandleFunc("POST /api/polka/webhooks", NewHandler(apiCfg.PostPolka))

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func HealthCheck(w http.ResponseWriter, request *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type CustomHandler func(w http.ResponseWriter, request *http.Request) error

func NewHandler(customHandler CustomHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := customHandler(w, r)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			respondWithError(w, err)
		}
	}
}
