package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajaen4/go-standard-lib-api/internal/db"
	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
	"github.com/ajaen4/go-standard-lib-api/pkg/encryption"
)

type PostChirpReq struct {
	Body string `json:"body"`
}

func (chirpReq *PostChirpReq) validate(r *http.Request) *api_errors.ClientErr {
	err := json.NewDecoder(r.Body).Decode(chirpReq)
	if err != nil {
		return &api_errors.ClientErr{
			HttpCode: http.StatusBadRequest,
			Message:  "Invalid JSON",
		}
	}

	apiErr := &api_errors.ClientErr{
		HttpCode: http.StatusBadRequest,
		Message:  "Invalid body parameters",
		Errors:   map[string]string{},
	}
	if len(chirpReq.Body) == 0 || len(chirpReq.Body) > 140 {
		apiErr.Errors["body"] = "invalid body"
	}

	if len(apiErr.Errors) > 0 {
		return apiErr
	}

	return nil
}

type GetChirpsReq struct {
	authorId int
	sortBy   string
}

func (req *GetChirpsReq) validate(r *http.Request) *api_errors.ClientErr {
	apiErr := &api_errors.ClientErr{
		HttpCode: http.StatusBadRequest,
		Message:  "invalid request params",
		Errors:   map[string]string{},
	}

	authorId := r.URL.Query().Get("author_id")
	if authorId != "" {
		authorIdInt, err := strconv.Atoi(authorId)
		if err != nil {
			apiErr.Errors["author_id"] = "invalid author_id query parameter"
		} else {
			req.authorId = authorIdInt
		}
	}

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "asc"
	}
	if sortBy != "asc" && sortBy != "desc" {
		apiErr.Errors["sort"] = "invalid sort query parameter"
	} else {
		req.sortBy = sortBy
	}

	if len(apiErr.Errors) > 0 {
		return apiErr
	}
	return nil
}

type ChirpReq struct {
	chirpID int
}

func (req *ChirpReq) validate(r *http.Request) *api_errors.ClientErr {
	apiErr := &api_errors.ClientErr{
		HttpCode: http.StatusBadRequest,
		Message:  "invalid request params",
		Errors:   map[string]string{},
	}

	reqChirpID := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(reqChirpID)
	if reqChirpID == "" || err != nil {
		apiErr.Errors["chirpID"] = "ChirpID not provided or invalid"
	} else {
		req.chirpID = chirpID
	}

	if len(apiErr.Errors) > 0 {
		return apiErr
	}
	return nil
}

func (apiCfg *ApiConfig) PostChirp(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := encryption.ValidateToken(tokenStr, apiCfg.JwtSecret)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	chirpReq := &PostChirpReq{}
	err = chirpReq.validate(r)
	if err != nil {
		return err
	}

	cleanWords := ProcessWords(chirpReq.Body)
	chirp, err := apiCfg.DB.CreateChirp(cleanWords, userId)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusCreated, chirp)
	return nil
}

func ProcessWords(str string) string {
	words := strings.Split(str, " ")
	newWords := make([]string, len(words))
	for i, word := range words {
		lWord := strings.ToLower(word)
		if lWord == "kerfuffle" || lWord == "sharbert" || lWord == "fornax" {
			newWords[i] = "****"
		} else {
			newWords[i] = word
		}
	}
	return strings.Join(newWords, " ")
}

func (apiCfg *ApiConfig) GetChirps(w http.ResponseWriter, request *http.Request) error {
	chirpsReq := GetChirpsReq{}
	chirpsReq.validate(request)

	var chirps []db.Chirp
	var err error
	if chirpsReq.authorId == 0 {
		chirps, err = apiCfg.DB.GetChirps(chirpsReq.sortBy)
	} else {
		chirps, err = apiCfg.DB.GetChirpsByAuthId(chirpsReq.authorId, chirpsReq.sortBy)
	}
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusOK, chirps)
	return nil
}

func (apiCfg *ApiConfig) GetChirp(w http.ResponseWriter, request *http.Request) error {
	chirpReq := ChirpReq{}
	clientErr := chirpReq.validate(request)
	if clientErr != nil {
		return clientErr
	}

	chirp, err := apiCfg.DB.GetChirp(chirpReq.chirpID)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusOK, chirp)
	return nil
}

func (apiCfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := encryption.ValidateToken(tokenStr, apiCfg.JwtSecret)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	chirpReq := ChirpReq{}
	err = chirpReq.validate(r)
	if err != nil {
		return err
	}

	err = apiCfg.DB.DeleteChirp(userId, chirpReq.chirpID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
