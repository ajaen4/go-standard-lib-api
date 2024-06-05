package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
)

type PolkaReq struct {
	Event string         `json:"event"`
	Data  map[string]int `json:"data"`
}

func (polkaReq *PolkaReq) validate(request *http.Request) error {
	err := json.NewDecoder(request.Body).Decode(polkaReq)
	if err != nil {
		return &api_errors.ClientErr{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid JSON",
		}
	}

	_, ok := polkaReq.Data["user_id"]
	if !ok {
		return &api_errors.ClientErr{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid body parameters",
			Errors:   map[string]string{"user_id": "invalid or not present user_id"},
		}
	}

	return nil
}

func (apiCfg *ApiConfig) PostPolka(w http.ResponseWriter, request *http.Request) error {
	authHeader := request.Header.Get("Authorization")
	apiKey := strings.Replace(authHeader, "ApiKey ", "", 1)
	if apiKey != apiCfg.PolkaKey {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = "Invalid ApiKey"
		return &apiErr
	}

	polkaReq := PolkaReq{}
	err := polkaReq.validate(request)
	if err != nil {
		return err
	}

	if polkaReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	err = apiCfg.DB.UserChirpyRed(polkaReq.Data["user_id"])
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
