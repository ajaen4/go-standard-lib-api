package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ajaen4/go-standard-lib-api/pkg/handlers"
)

func TestPostUser(t *testing.T) {
	apiCfg := setupTestCfg(t)

	userReq := &handlers.UserReq{
		Email:    "test@email.com",
		Password: "testPassword",
	}
	jsonReq, err := json.Marshal(userReq)
	if err != nil {
		t.Fatal(err)
	}

	readerReq := strings.NewReader(string(jsonReq))
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/api/users", readerReq)
	if err != nil {
		t.Fatal(err)
	}

	err1 := apiCfg.PostUser(w, r)
	if err != nil {
		t.Fatal(err1)
	}

	response := handlers.UserResp{}
	expected := handlers.UserResp{
		Id:          1,
		Email:       userReq.Email,
		IsChirpyRed: false,
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response != expected {
		t.Errorf(
			"handler returned wrong response: got %v want %v",
			response,
			expected,
		)
	}

	tearDownTestCfg(t, apiCfg)
}
