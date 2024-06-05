package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
)

type errorResponse struct {
	Code    int    `json:"int"`
	Message string `json:"error"`
}

func respondWithError(w http.ResponseWriter, err error) {
	if clientErr, ok := err.(*api_errors.ClientErr); ok {
		respondWithJSON(w, clientErr.HttpCode, clientErr)
	} else {
		respondWithJSON(w, http.StatusInternalServerError, errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	jsonPay, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error when marshaling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(jsonPay)
}
