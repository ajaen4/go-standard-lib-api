package handlers

import (
	"fmt"
	"net/http"
)

func (apiCfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.FileserverHits++
		http.StripPrefix("/app", next).ServeHTTP(w, r)
	})
}

func (apiCfg *ApiConfig) MetricsReset(w http.ResponseWriter, r *http.Request) {
	apiCfg.FileserverHits = 0
}

func (apiCfg *ApiConfig) MetricsCount(w http.ResponseWriter, request *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	
	</html>
	`, apiCfg.FileserverHits)))
}
