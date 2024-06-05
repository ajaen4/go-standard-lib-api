package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type InternalErrResp struct {
	Code    int    `json:"int"`
	Message string `json:"error"`
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
