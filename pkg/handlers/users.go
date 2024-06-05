package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajaen4/go-standard-lib-api/pkg/api_errors"
	"github.com/ajaen4/go-standard-lib-api/pkg/encryption"
)

type UserReq struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

func (userReq *UserReq) validate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(userReq)

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
	if len(userReq.Email) == 0 {
		apiErr.Errors["email"] = "invalid email"
	}
	if len(userReq.Password) == 0 {
		apiErr.Errors["password"] = "invalid email"
	}

	if len(apiErr.Errors) > 0 {
		return apiErr
	}

	return nil
}

type UserResp struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type LogInResp struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
	Token       string `json:"token"`
	RefToken    string `json:"refresh_token"`
}

type TokenResp struct {
	Token string `json:"token"`
}

func (apiCfg *ApiConfig) PostUser(w http.ResponseWriter, request *http.Request) error {
	userReq := &UserReq{}
	if reqErr := userReq.validate(request); reqErr != nil {
		return reqErr
	}

	User, err := apiCfg.DB.CreateUser(userReq.Email, userReq.Password)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusCreated, UserResp{
		Id:          User.Id,
		Email:       User.Email,
		IsChirpyRed: User.IsChirpyRed,
	})
	return nil
}

func (apiCfg *ApiConfig) PutUser(w http.ResponseWriter, request *http.Request) error {

	authHeader := request.Header.Get("Authorization")
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

	id, err := strconv.Atoi(subject)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	userReq := &UserReq{}
	if reqErr := userReq.validate(request); reqErr != nil {
		return reqErr
	}

	user, err := apiCfg.DB.UpdateUser(id, userReq.Email, userReq.Password)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusOK, UserResp{
		Email:       user.Email,
		Id:          id,
		IsChirpyRed: user.IsChirpyRed,
	})
	return nil
}

func (apiCfg *ApiConfig) PostLogin(w http.ResponseWriter, request *http.Request) error {
	userReq := &UserReq{}
	if reqErr := userReq.validate(request); reqErr != nil {
		return reqErr
	}

	base64RefToken, err := encryption.CreateRefToken()
	if err != nil {
		return err
	}

	User, err := apiCfg.DB.Login(userReq.Email, userReq.Password)
	if err != nil {
		return err
	}

	signedToken, err := encryption.CreateToken(User.Id, apiCfg.JwtSecret)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	err = apiCfg.DB.SaveRefToken(User.Id, base64RefToken)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusOK, LogInResp{
		Id:          User.Id,
		Email:       User.Email,
		IsChirpyRed: User.IsChirpyRed,
		Token:       signedToken,
		RefToken:    base64RefToken,
	})
	return nil
}

func (apiCfg *ApiConfig) PostRefToken(w http.ResponseWriter, request *http.Request) error {
	authHeader := request.Header.Get("Authorization")
	refreshToken := strings.Replace(authHeader, "Bearer ", "", 1)

	User, err := apiCfg.DB.ValidateRefToken(refreshToken)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	signedToken, err := encryption.CreateToken(User.Id, apiCfg.JwtSecret)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusOK, TokenResp{
		Token: signedToken,
	})
	return nil
}

func (apiCfg *ApiConfig) PostRevokeToken(w http.ResponseWriter, request *http.Request) error {
	authHeader := request.Header.Get("Authorization")
	refreshToken := strings.Replace(authHeader, "Bearer ", "", 1)

	err := apiCfg.DB.RevokeRefToken(refreshToken)
	if err != nil {
		apiErr := api_errors.UnauthErr
		apiErr.LogMess = err.Error()
		return &apiErr
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
